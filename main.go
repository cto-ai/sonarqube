package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"net/http"
	"net/url"

	"os"
	"os/exec"

	ctoai "github.com/cto-ai/sdk-go"
	"github.com/sethvargo/go-diceware/diceware"
	"gopkg.in/src-d/go-git.v4"
)

const cloneDir = "/ops/proj"

var client = ctoai.NewClient()

//badErrHandler prints and panics on error
func badErrHandler(err error) {
	if err != nil {
		client.Ux.Print(fmt.Sprint(err))
		panic(err)
	}
}

//printWrapper fixes client.Ux.Print, allowing it to print arbitrary, multiple arguments ala fmt.Println
func printWrapper(a ...interface{}) {
	client.Ux.Print(fmt.Sprint(a...))
}

//Ask the user for which repo to analyse
func promptRepo() (string, error) {
	gitURL, err := client.Prompt.Input("gitproject", "Please enter your publicly accessible, http(s) git repository URL to analyse", ctoai.OptInputAllowEmpty(false), ctoai.OptInputFlag("g"))
	if err != nil {
		printWrapper("❌ Encountered an error: ", err)
		return promptRepo()
	}
	if gitURL == "" {
		printWrapper("❌ URL cannot be empty!")
		return promptRepo()
	}
	//check if a valid URL
	u, err := url.ParseRequestURI(gitURL)
	if err != nil {
		printWrapper("❌ Invalid URL detected ", err)
		return promptRepo()
	}
	//TODO: check that it's an actual valid Git repository
	if !strings.HasPrefix(u.Scheme, "http") {
		printWrapper("❌ We only support cloning with http(s), we do not support ", u.Scheme)
		return promptRepo()
	}
	return gitURL, err
}

//validhost is used to avoid asking user for host again if they give an invalid token
//pass in "" for the first time calling this function
func promptCredentials(validhost string) (host string, token string) {
	var err error = nil
	if validhost == "" {
		host, err = client.Sdk.GetSecret("sonarHost")
		if err != nil {
			printWrapper("❌ That URL does not appear to be a valid SonarQube server, please try again!")
			return promptCredentials("")
		}
		host = strings.TrimSuffix(host, "/")
		err = ValidateHost(host)
		if err != nil {
			printWrapper("❌ That URL does not appear to be a valid SonarQube server, please try again!")
			return promptCredentials("")
		}
	} else {
		host = validhost
	}

	token, err = client.Sdk.GetSecret("sonarToken")
	if err != nil {
		printWrapper("❌ That token does not appear to be a valid SonarQube administrative token, please try again!")
		return promptCredentials(host)
	}
	err = ValidateToken(host, token)
	if err != nil {
		printWrapper("❌ That token does not appear to be a valid SonarQube administrative token, please try again!")
		return promptCredentials(host)
	}
	return
}

//printBig chunks out the input string
func printBig(input string) {
	const chunkSize = 2000
	chunks := len(input) / chunkSize
	for index := 0; index < chunks; index++ {
		printWrapper(input[index*chunkSize : (index+1)*chunkSize])
	}
	printWrapper(input[chunks*chunkSize : (chunks*chunkSize)+(len(input)%chunkSize)])
}

//takes in the output from SonarScanner
//and breaks it into chunks to evade Slack limits
func parseAndPrintSSOutput(output []byte) {
	const numberOfOutputLines = 12
	ssOutputArray := strings.Split(string(output), "\n")
	if len(ssOutputArray) < numberOfOutputLines {
		printBig(string(output))
	}
	tail := ssOutputArray[len(ssOutputArray)-numberOfOutputLines:]
	finalOutput := make([]string, len(tail))

	const separator = ": "
	for index, line := range tail {
		colIndex := strings.Index(line, separator)
		if colIndex == -1 {
			finalOutput[index] = line
		} else {
			finalOutput[index] = line[colIndex+len(separator):]
		}
	}
	printBig(strings.Join(finalOutput, "\n"))
}

//simply clone to cloneDir
func cloneRepo(gitURL string, cloneChan chan<- error) {
	_, err := git.PlainClone(cloneDir, false, &git.CloneOptions{
		URL:          gitURL,
		Depth:        1,
		SingleBranch: true,
	})
	cloneChan <- err
}

//ValidateHost checks if the given host is a valid SonarQube server
//returns nil if it is valid, or err otherwise
func ValidateHost(host string) (err error) {
	verResp, err := http.Get(host + "/api/server/version")
	if err != nil {
		return err
	}
	if verResp.StatusCode != 200 {
		return fmt.Errorf("❌ Invalid host, unable to obtain server version, status code %d", verResp.StatusCode)
	}
	return nil
}

//ValidateToken checks if the given token is valid
//Assumes that host is a valid SonarQube server
//returns nil if it is valid, or err otherwise
func ValidateToken(host string, token string) (err error) {
	req, err := http.NewRequest("GET", host+"/api/system/ping", nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(token, "")
	pingResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer pingResp.Body.Close()
	body, err := ioutil.ReadAll(pingResp.Body)
	if err != nil {
		return err
	}

	if pingResp.StatusCode != 200 {
		return fmt.Errorf("❌ Invalid token, status code expected 200 got %d", pingResp.StatusCode)
	} else if string(body) != "pong" {
		maxlen := len(body)
		if len(body) > 32 {
			maxlen = 32
		}
		return fmt.Errorf("❌ Invalid host, response expected pong got \n  %s", string(body[0:maxlen]))
	} else {
		return nil
	}
}

//represents the credentials for the users to view the final report
type ssUser struct {
	login    string
	password string
}

//creates a temp user with default permissions
//the SS server should be configured to give no permissions by default
//preventing users from seeing any project other than the ones they get explicit permissions for
//eg the project they ran the analysis on
func createTempUser(host string, adminToken string) (user ssUser, err error) {
	//double duty as unique name/login
	suffix, err := diceware.Generate(2)
	if err != nil {
		suffix = []string{string(time.Now().Nanosecond()), string(time.Now().UnixNano())}
	}
	password, err := diceware.Generate(3)
	if err != nil {
		password = []string{string(time.Now().Nanosecond()), os.Getenv("OPS_TEAM_NAME"), string(time.Now().UnixNano())}
	}
	err = nil
	user.login = os.Getenv("OPS_TEAM_NAME") + "-" + strings.Join(suffix, "")
	user.password = strings.Join(password, "")

	createPostForm := url.Values{}
	createPostForm.Set("login", user.login)
	createPostForm.Set("name", user.login)
	createPostForm.Set("password", user.password)

	req, err := http.NewRequest("POST", host+"/api/users/create", strings.NewReader(createPostForm.Encode()))
	if err != nil {
		return
	}
	req.SetBasicAuth(adminToken, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	createResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer createResp.Body.Close()
	if createResp.StatusCode != 200 {
		err = errors.New("Failed to create user, expected status code 200 got" + string(createResp.StatusCode))
	}
	return
}

//give user the permissions read access for the specified project
func grantUserPerms(host string, adminToken string, user ssUser, projectKey string) (err error) {
	//scan left out, only those with access to the admin token can run analysis?
	// "scan",
	necessaryPermissions := []string{"codeviewer", "user"}

	errChan := make(chan error, len(necessaryPermissions))

	requestOnePerm := func(permission string, errChan chan<- error) {
		grantPermForm := url.Values{}
		grantPermForm.Set("login", user.login)
		grantPermForm.Set("permission", permission)
		grantPermForm.Set("projectKey", projectKey)

		req, err := http.NewRequest("POST", host+"/api/permissions/add_user", strings.NewReader(grantPermForm.Encode()))
		if err != nil {
			errChan <- err
			return
		}
		req.SetBasicAuth(adminToken, "")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		createResp, err := http.DefaultClient.Do(req)
		if err != nil {
			errChan <- err
			return
		}

		createResp.Body.Close()
		if createResp.StatusCode != 204 {
			errChan <- errors.New(fmt.Sprint("Failed to grant ", user.login, "permission. Received status code ", string(createResp.StatusCode)))
			return
		}
		errChan <- nil
	}

	for _, onePerm := range necessaryPermissions {
		go requestOnePerm(onePerm, errChan)
	}

	for range necessaryPermissions {
		if err = <-errChan; err != nil {
			printWrapper(err)
			return err
		}
	}
	return nil
}

func main() {
    printLogo(client)

	cloneChan := make(chan error, 1)

	//parse arguments as repo if available, else prompt
	var curRepo string
	var err error
	repo := os.Args[1:]
	if len(repo) >= 1 {
		//TODO: validate like the prompt
		curRepo = repo[0]
	} else {
		curRepo, err = promptRepo()
		badErrHandler(err)
	}

	//run the clone in the background while the prompts buy time
	go cloneRepo(curRepo, cloneChan)
	host, token := promptCredentials("")
	user, err := createTempUser(host, token)
	badErrHandler(err)

	client.Ux.SpinnerStart("📠 Cloning your Git repository...")
	//if the repo is small, this will happen instantly like magic
	badErrHandler(<-cloneChan)
	client.Ux.SpinnerStop("♊ Cloning finished!")

	//can't access user, so team name will have to do
	projectKey := os.Getenv("OPS_TEAM_NAME") + "-" + curRepo
	ssCommand := exec.Command("sonar-scanner", "-Dsonar.host.url="+host, "-Dsonar.login="+token, "-Dsonar.projectKey="+projectKey)
	ssCommand.Dir = cloneDir

	client.Ux.SpinnerStart("🔍 Running analysis...")
	ssOutput, err := ssCommand.CombinedOutput()
	if err != nil {
		parseAndPrintSSOutput(ssOutput)
		printBig(err.Error())
		panic(err)
	}
	grantPermErr := grantUserPerms(host, token, user, projectKey)
	client.Ux.SpinnerStop("✅ Analysis finished!")

	parseAndPrintSSOutput(ssOutput)
	if grantPermErr != nil {
		printWrapper("Error granting your new user the permissions to access the analysis!", grantPermErr)
	} else {
		printWrapper("Access your analysis with the credentials `", user.login, "` and password `", user.password+"`")
	}
}
