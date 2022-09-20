package main

import (
	"bufio"
	_ "context"
	"encoding/json"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	_ "log"
	"net"
	"net/url"
	"os"
	_ "time"
)

func main() {

	fmt.Println("Thanks for using the Auto-Backup tool! This tool will automatically backup your files to a remote server.")
	fmt.Println("Trying to read config file [config.json]...")

	// Get SFTP To Go URL from environment
	rawurl := os.Getenv("SFTPTOGO_URL")

	parsedUrl, err := url.Parse(rawurl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse SFTP To Go URL: %s\n", err)
		os.Exit(1)
	}

	// Get user name and pass
	user := parsedUrl.User.Username()
	pass, _ := parsedUrl.User.Password()

	// Parse Host and Port
	host := parsedUrl.Host
	// Default SFTP port
	port := 22

	hostKey := getHostKey(host)

	fmt.Fprintf(os.Stdout, "Connecting to %s ...\n", host)

	var auths []ssh.AuthMethod

	// Try to use $SSH_AUTH_SOCK which contains the path of the unix file socket that the sshd agent uses
	// for communication with other processes.
	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
	}

	// Use password authentication if provided
	if pass != "" {
		auths = append(auths, ssh.Password(pass))
	}

	// Initialize client configuration
	remoteConfig := ssh.ClientConfig{
		User: user,
		Auth: auths,
		// Uncomment to ignore host key check
		//HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	// Connect to server
	conn, err := ssh.Dial("tcp", addr, &remoteConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connecto to [%s]: %v\n", addr, err)
		os.Exit(1)
	}

	defer conn.Close()

	// Create new SFTP client
	sc, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start SFTP subsystem: %v\n", err)
		os.Exit(1)
	}
	defer sc.Close()

	test()

	jsonFile, err := os.Open("users.json")
	if err != nil {
		// handle error
		// if the file doesn't exist, create it
		fmt.Println("Config file not found. Creating new config file...")
		// create a new file
		newFile, err := os.Create("users.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		// write to the file with the Config struct with pretty print
		prettyJSON, err := json.MarshalIndent(Config{
			RemotePath: "",
			LocalPaths: []string{"/Users/max/Downloads", "/Users/max/Downloads"},
			Delay:      4,
		}, "", "    ")
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = newFile.Write(prettyJSON)
	}
	fmt.Println("Successfully opened users.json")
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(jsonFile)

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Println(err)
		return
	}

	ZipWriter(".idea/", "test")
	for i := 0; i < len(config.LocalPaths); i++ {
		fmt.Println(config.LocalPaths[i])
	}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		fmt.Println("> Echo:", input.Text())
	}

}

type Config struct {
	IP         string   `json:"remote_ip"`
	Password   string   `json:"remote_password"`
	Port       int      `json:"remote_port"`
	RemotePath string   `json:"remote_save_path"`
	LocalPaths []string `json:"local_path"`
	Delay      int      `json:"delay"`
}
