package cmd

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/ystyle/jvms/utils/file"
	"github.com/ystyle/jvms/utils/jdk"
	"github.com/ystyle/jvms/utils/web"
	"os"
	"os/exec"
	"path/filepath"
)

func initAction(c *cli.Context) error {
	if c.IsSet("jvms_home") {
		AppConfig.JvmsHome = c.String("jvms_home")
	}
	if c.IsSet("java_home") || AppConfig.JavaHome == "" {
		AppConfig.JavaHome = c.String("java_home")
	}
	cmd := exec.Command("cmd", "/C", "setx", "JAVA_HOME", AppConfig.JavaHome, "/M")
	err := cmd.Run()
	if err != nil {
		return errors.New("set Environment variable `JAVA_HOME` failure: Please run as admin user")
	}
	fmt.Println("set `JAVA_HOME` Environment variable to ", AppConfig.JavaHome)

	if c.IsSet("originalpath") || AppConfig.Originalpath == "" {
		AppConfig.Originalpath = c.String("originalpath")
	}
	path := fmt.Sprintf(`%s%sbin;%s;%s`, AppConfig.JavaHome, filepath.Separator, os.Getenv("PATH"), file.GetCurrentPath())
	cmd = exec.Command("cmd", "/C", "setx", "path", path, "/m")
	err = cmd.Run()
	if err != nil {
		return errors.New("set Environment variable `PATH` failure: Please run as admin user")
	}
	fmt.Println("add jvms.exe to `path` Environment variable")
	return nil
}

func configAction(c *cli.Context) error {
	fmt.Println(ShowConfig())
	return nil
}

func listAction(c *cli.Context) error {
	fmt.Println("Installed jdk (* marks in use):")
	v := jdk.GetInstalled(AppConfig.store)
	for i, version := range v {
		str := ""
		if AppConfig.CurrentJDKVersion == version {
			str = fmt.Sprintf("%s  * %d) %s", str, i+1, version)
		} else {
			str = fmt.Sprintf("%s    %d) %s", str, i+1, version)
		}
		fmt.Printf(str + "\n")
	}
	if len(v) == 0 {
		fmt.Println("No installations recognized.")
	}
	return nil
}

func installAction(c *cli.Context) error {
	if AppConfig.Proxy != "" {
		web.SetProxy(AppConfig.Proxy)
	}
	v := c.Args().Get(0)
	if v == "" {
		return errors.New("invalid version., Type \"jvms rls\" to see what is available for install")
	}

	if jdk.IsVersionInstalled(AppConfig.store, v) {
		fmt.Println("Version " + v + " is already installed.")
		return nil
	}

	if !file.Exists(AppConfig.download) {
		err := os.MkdirAll(AppConfig.download, 0777)
		if err != nil {
			return err
		}
	}
	if !file.Exists(AppConfig.store) {
		err := os.MkdirAll(AppConfig.store, 0777)
		if err != nil {
			return err
		}
	}

	versions := CacheGetVersion()
	if len(versions) == 0 {
		var err error
		versions, err = jdk.RemoteJdkVersions()
		if err != nil {
			return err
		}
	}

	installed := file.Exists(filepath.Join(AppConfig.store, v))
	if installed {
		fmt.Printf("JDK %v ")
	}

	for _, version := range versions {
		if version.Version == v {
			dlzipfile, success := web.GetJDK(AppConfig.download, v, version.Url)
			if success {
				// Extract jdk to the temp directory
				jdktempfile := filepath.Join(AppConfig.download, fmt.Sprintf("%s_temp", v))
				if file.Exists(jdktempfile) && c.BoolT("f") {
					err := os.RemoveAll(jdktempfile)
					if err != nil {
						panic(err)
					}
				}

				err := file.Unzip(dlzipfile, jdktempfile)
				if err != nil {
					return fmt.Errorf("unzip failed: %w", err)
				}

				// Copy the jdk files to the installation directory
				temJavaHome := jdk.GetJavaHome(jdktempfile)
				err = os.Rename(temJavaHome, filepath.Join(AppConfig.store, v))
				if err != nil {
					return fmt.Errorf("unzip failed: %w", err)
				}

				// Remove the temp directory
				// may consider keep the temp files here
				_ = os.RemoveAll(jdktempfile)

				fmt.Println("Installation complete. If you want to use this version, type\n\njvms switch", v)
			} else {
				fmt.Println("Could not download JDK " + v + " executable.")
			}
			return nil
		}
	}
	return errors.New("invalid version., Type \"jvms rls\" to see what is available for install")
}

func addAction(c *cli.Context) error {
	version := c.String("version")
	if version == "" {
		return errors.New("you should input a version name")
	}
	path := c.Args().Get(0)
	if path == "" || !filepath.IsAbs(path) {
		return errors.New("you should input a Local JDK path")
	}

	if jdk.GetJavaHome(path) == "" {
		return errors.New("args value " + path + " is not JDK")
	}

	if jdk.IsVersionInstalled(AppConfig.store, version) {
		fmt.Printf("jdk %s is already installed. ", version)
		return nil
	}
	err := os.Symlink(path, filepath.Join(AppConfig.store, version))
	if err != nil {
		return errors.New("Add local jdk failed, " + err.Error())
	}
	fmt.Printf("Add local jdk version=%s path=%s complete.\n", version, path)
	return nil
}

func switchAction(c *cli.Context) error {
	v := c.Args().Get(0)
	if v == "" {
		return errors.New("you should input a version, Type \"jvms list\" to see what is installed")
	}
	if !jdk.IsVersionInstalled(AppConfig.store, v) {
		fmt.Printf("jdk %s is not installed. ", v)
		return nil
	}
	// Create or update the symlink
	if file.Exists(AppConfig.JavaHome) {
		err := os.Remove(AppConfig.JavaHome)
		if err != nil {
			return errors.New("Switch jdk failed, please manually remove " + AppConfig.JavaHome)
		}
	}
	cmd := exec.Command("cmd", "/C", "setx", "JAVA_HOME", AppConfig.JavaHome, "/M")
	err := cmd.Run()
	if err != nil {
		return errors.New("set Environment variable `JAVA_HOME` failure: Please run as admin user")
	}
	err = os.Symlink(filepath.Join(AppConfig.store, v), AppConfig.JavaHome)
	if err != nil {
		return errors.New("Switch jdk failed, " + err.Error())
	}
	fmt.Println("Switch success.\nNow using JDK " + v)
	AppConfig.CurrentJDKVersion = v
	return nil
}

func removeAction(c *cli.Context) error {
	v := c.Args().Get(0)
	if v == "" {
		return errors.New("you should input a version, Type \"jvms list\" to see what is installed")
	}
	if jdk.IsVersionInstalled(AppConfig.store, v) {
		fmt.Printf("Remove JDK %s ...\n", v)
		if AppConfig.CurrentJDKVersion == v {
			os.Remove(AppConfig.JavaHome)
		}
		dir := filepath.Join(AppConfig.store, v)
		e := os.RemoveAll(dir)
		if e != nil {
			fmt.Println("Error removing jdk " + v)
			fmt.Println("Manually remove " + dir + ".")
		} else {
			fmt.Printf(" done")
		}
	} else {
		fmt.Println("jdk " + v + " is not installed. Type \"jvms list\" to see what is installed.")
	}
	return nil
}

func originsAction(c *cli.Context) error {
	fmt.Printf("\nFor a complete list, visit \n")
	for _, source := range jdk.JdkOrigins {
		fmt.Printf("\n %s: %s\n", source.OriginName(), source.OriginDesc())
	}
	return nil
}

func versionsAction(c *cli.Context) error {
	if AppConfig.Proxy != "" {
		web.SetProxy(AppConfig.Proxy)
	}

	var versions = CacheGetVersion()
	var err error
	fmt.Println(versions)
	if len(versions) == 0 || c.Bool("force") {
		versions, err = jdk.RemoteJdkVersions()
		if err != nil {
			return err
		}
		if c.Bool("cache") {
			CachePutVersion(versions)
		}
	}
	for i, version := range versions {
		fmt.Printf("    %d) %s\n", i+1, version.Version)
		if !c.Bool("a") && i >= 9 {
			fmt.Println("\nuse \"jvm rls -a\" show all the versions ")
			break
		}
	}
	if len(versions) == 0 {
		fmt.Println("No availabled jdk veriosn for download.")
	}
	return nil
}

func proxyAction(c *cli.Context) error {
	if c.Bool("show") {
		fmt.Printf("Current proxy: %s\n", AppConfig.Proxy)
		return nil
	}
	if c.IsSet("set") {
		AppConfig.Proxy = c.String("set")
	}
	return nil
}
