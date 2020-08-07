package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/mervick/aes-everywhere/go/aes256"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const websiteUrl = "https://example.com/idigo"

func stopServices() {
	servicesToStop := []string{"HTTPDebuggerSvc.exe", "HTTPDebuggerUI.exe", "Proxifier.exe", "Fiddler4.exe", "Wireshark.exe"}
	for i := 0; i < len(servicesToStop); i++ {
		cmd := exec.Command("cmd.exe", "/C", `taskkill /im `+servicesToStop[i]+` /f`)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	cmd := exec.Command("cmd.exe", "/C", `net stop npf`)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

const TH32CS_SNAPPROCESS = 0x00000002

type WindowsProcess struct {
	ProcessID       int
	ParentProcessID int
	Exe             string
}

func processes() ([]WindowsProcess, error) {
	handle, err := windows.CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(handle)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	// get the first process
	err = windows.Process32First(handle, &entry)
	if err != nil {
		return nil, err
	}

	results := make([]WindowsProcess, 0, 50)
	for {
		results = append(results, newWindowsProcess(&entry))

		err = windows.Process32Next(handle, &entry)
		if err != nil {
			// windows sends ERROR_NO_MORE_FILES on last process
			if err == syscall.ERROR_NO_MORE_FILES {
				return results, nil
			}
			return nil, err
		}
	}
}

func findProcessByName(processes []WindowsProcess, name string) *WindowsProcess {
	for _, p := range processes {
		if strings.ToLower(p.Exe) == strings.ToLower(name) {
			return &p
		}
	}
	return nil
}

func newWindowsProcess(e *windows.ProcessEntry32) WindowsProcess {
	// Find when the string ends for decoding
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}

	return WindowsProcess{
		ProcessID:       int(e.ProcessID),
		ParentProcessID: int(e.ParentProcessID),
		Exe:             syscall.UTF16ToString(e.ExeFile[:end]),
	}
}
func DownloadFile(url string, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
func MessageBox(hwnd uintptr, caption, title string, flags uint) int {
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(flags))

	return int(ret)
}
func scan(javawStrings []string, lsassStrings []string, MsMpEngStrings []string, explorerStrings []string, dwmStrings []string, w *astilectron.Window, pin string, lastModifiedRB int64) <-chan string {
	r := make(chan string)

	go func() {
		defer close(r)

		DownloadFile(websiteUrl+"/ms-files/64.exe", os.Getenv("TEMP")+"\\idigo64.exe")
		procs, err := processes()
		if err != nil {
			log.Fatal(err)
		}

		var foundCheats string
		javaw := findProcessByName(procs, "javaw.exe")
		if javaw != nil {
			cmd := exec.Command("cmd.exe", "/C", os.Getenv("TEMP")+"\\idigo64.exe "+"-nh -pid "+strconv.Itoa(javaw.ProcessID))
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

			javaw, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatal(err)
			}
			javawSstr := fmt.Sprintf("%s", javaw)
			for i := 0; i < len(javawStrings); i++ {
				if strings.Contains(javawSstr, strings.Split(string(javawStrings[i]), ":::")[0]) {
					index := strings.Index(javawSstr, strings.Split(string(javawStrings[i]), ":::")[0])
					// minus := 1
					// if index == 0 {
					// 	minus--
					// }
					// start := index - minus
					// end := index + len(strings.Split(string(javawStrings[i]), ":::")[0]) + 1

					// firstChar := javawSstr[start : start+1]
					// secondChar := javawSstr[end-minus : end]
					// if len(javawSstr) == end-1 {
					// 	secondChar = string([]byte{0xfd})
					// }
					// if index == 0 {
					// 	firstChar = string([]byte{0xfd})
					// }
					// firstChar = strings.ReplaceAll(firstChar, "\n", string([]byte{0xfd}))
					// secondChar = strings.ReplaceAll(secondChar, "\n", string([]byte{0xfd}))
					//if !utf8.Valid([]byte(firstChar)) && !utf8.Valid([]byte(secondChar)) {
					start := index
					end := index + len(strings.Split(string(javawStrings[i]), ":::")[0])
					if !strings.Contains(javawSstr[start:end+20], ".zip") && !strings.Contains(javawSstr[start:end+20], ".rar") {
						fmt.Println("Found: " + strings.Split(string(javawStrings[i]), ":::")[0])
						foundCheats += strings.Split(string(javawStrings[i]), ":::")[1] + "|||"
					}
					//}
				}
			}
		}
		explorer := findProcessByName(procs, "explorer.exe")
		if explorer != nil {
			cmd := exec.Command("cmd.exe", "/C", os.Getenv("TEMP")+"\\idigo64.exe "+"-nh -pid "+strconv.Itoa(explorer.ProcessID))
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

			stringsFO, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatal(err)
			}
			actualStrings := fmt.Sprintf("%s", stringsFO)
			for i := 0; i < len(explorerStrings); i++ {
				if strings.Contains(actualStrings, strings.Split(string(explorerStrings[i]), ":::")[0]) {
					fmt.Println("Found: " + strings.Split(string(explorerStrings[i]), ":::")[0])
					foundCheats += strings.Split(string(explorerStrings[i]), ":::")[1] + "|||"
				}
			}
		}
		msmpeng := findProcessByName(procs, "MsMpEng.exe")
		if msmpeng != nil {
			cmd := exec.Command("cmd.exe", "/C", os.Getenv("TEMP")+"\\idigo64.exe "+"-nh -pid "+strconv.Itoa(msmpeng.ProcessID))
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

			stringsFO, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatal(err)
			}
			actualStrings := fmt.Sprintf("%s", stringsFO)
			for i := 0; i < len(MsMpEngStrings); i++ {
				if strings.Contains(actualStrings, strings.Split(string(MsMpEngStrings[i]), ":::")[0]) {
					fmt.Println("Found: " + strings.Split(string(MsMpEngStrings[i]), ":::")[0])
					foundCheats += strings.Split(string(MsMpEngStrings[i]), ":::")[1] + "|||"
				}
			}
		}
		lsass := findProcessByName(procs, "lsass.exe")
		if lsass != nil {
			cmd := exec.Command("cmd.exe", "/C", os.Getenv("TEMP")+"\\idigo64.exe "+"-nh -pid "+strconv.Itoa(lsass.ProcessID))
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

			stringsFO, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatal(err)
			}
			actualStrings := fmt.Sprintf("%s", stringsFO)
			for i := 0; i < len(explorerStrings); i++ {
				if strings.Contains(actualStrings, strings.Split(string(lsassStrings[i]), ":::")[0]) {
					fmt.Println("Found: " + strings.Split(string(lsassStrings[i]), ":::")[0])
					foundCheats += strings.Split(string(lsassStrings[i]), ":::")[1] + "|||"
				}
			}
		}
		_ = os.Remove(os.Getenv("TEMP") + "\\idigo64.exe")

		client := &http.Client{}
		req, err := http.NewRequest("GET", websiteUrl+"/sendResults.php", nil)
		req.Header.Add("User-Agent", `idigo (v0.1a): client`)
		req.Header.Add("ST", aes256.Encrypt("6mYiPz5u1WO62MM80kyZpBEQXHlD0Ho8nEwqzcegSj1ZeqqeplZHcRAry4e0lKWyUnn7VKiSOdrWR897BGtxaxLJuAJze96mza2", "4Oe7EmckEVKuogjcoLQWAaVhJAkIs6PT"))
		req.Header.Add("PIN", aes256.Encrypt(pin, "PINMY8n7aVQ7WX03aAHV8mbzFBBsEWIO"))

		if len(foundCheats) == 0 {
			foundCheats = "None"
		}
		user, _ := user.Current()

		k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
		if err != nil {
			fmt.Println("Error CurrentVersion not found")
		}
		defer k.Close()

		cv, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
		if err != nil {
			cv = 7
		}

		build, _, err := k.GetStringValue("CurrentBuild")
		if err != nil {
			build = "Not found"
		}

		encryptedResults := aes256.Encrypt("Client::IdigoValidRequest|||IDIGO|||"+foundCheats+"|||"+fmt.Sprintf("%s (%s)", user.Username, user.Name)+"|||"+fmt.Sprintf("%d", lastModifiedRB)+"|||"+fmt.Sprintf("Windows %d (%s)", uint64(cv), build), "Uv21t0oNjUHlHSiMcTO8cIsWYa77vo7y")
		req.Header.Add("RESULTS", encryptedResults)
		fmt.Println("Sending resulkts")
		stopServices()
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		bodytemp, err := ioutil.ReadAll(resp.Body)

		body := aes256.Decrypt(string(bodytemp), "Uv21t0oNjUHlHSiMcTO8cIsWYa77vo7y")
		fmt.Println("Body temp: " + string(bodytemp))
		fmt.Println("Body: " + string(body))
		resultsresponse := strings.Split(string(body), "|||")
		if resultsresponse[0] != "API::VALID_REQUEST" {
			w.Hide()
			MessageBox(0, "User might be using some kind of Idigo Killer\n\nCheats found: "+foundCheats, "Invalid Request", 0)
			w.Close()
			os.Exit(0)
		}

		w.SendMessage("Idigo::ScanFinished")
		w.ExecuteJavaScript(`document.location = websiteUrl + "/last.html?v=" + Math.random();`)
	}()
	return nil
}

func randstr(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func main() {
	lastModifiedRBTemp, _ := os.Stat("C:\\$Recycle.Bin")
	lastModifiedRB := lastModifiedRBTemp.ModTime().Unix()

	rand.Seed(time.Now().UnixNano())
	// Set logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	// Create astilectron
	a, err := astilectron.New(l, astilectron.Options{
		AppName: "Idigo",
	})
	if err != nil {
		l.Fatal(fmt.Errorf("main: creating astilectron failed: %w", err))
	}
	defer a.Close()

	// Handle signals
	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		l.Fatal(fmt.Errorf("main: starting astilectron failed: %w", err))
	}

	randomstring := randstr(16)
	procs, err := processes()
	if err != nil {
		log.Fatal(err)
	}

	stopServices()
	DownloadFile(websiteUrl+"/favicon.ico", os.Getenv("TEMP")+"\\"+randomstring+"idigo.ico")
	// New window
	var w *astilectron.Window
	url := websiteUrl + "/index.html?v=" + fmt.Sprintf("%d", rand.Int())

	tempJavawTest := findProcessByName(procs, "javaw.exe")
	if tempJavawTest == nil {
		url = websiteUrl + "/error?e=NOJAVAW"
	}

	if w, err = a.NewWindow(url, &astilectron.WindowOptions{
		Center:          astikit.BoolPtr(true),
		Height:          astikit.IntPtr(410),
		Width:           astikit.IntPtr(640),
		Frame:           astikit.BoolPtr(false),
		BackgroundColor: astikit.StrPtr("#121212"),
		AlwaysOnTop:     astikit.BoolPtr(true),
		Icon:            astikit.StrPtr(os.Getenv("TEMP") + "\\" + randomstring + "idigo.ico"),
		Resizable:       astikit.BoolPtr(false),
	}); err != nil {
		l.Fatal(fmt.Errorf("main: new window failed: %w", err))
	}

	// This will listen to messages sent by Javascript
	w.OnMessage(func(m *astilectron.EventMessage) interface{} {
		// Unmarshal
		var s string
		m.Unmarshal(&s)

		// Process message
		if s == "hello" {
			return "world"
		}

		split := strings.Split(s, "|")
		if strings.Contains(s, "Client::OpenWebsite") {
			cmd := exec.Command("cmd.exe", "/C", "start "+split[1])
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

			stdoutStderr, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(stdoutStderr)
		}
		if strings.Contains(s, "Client::Scan") {
			match, _ := regexp.MatchString("[^0-9]", split[1])
			if len(split[1]) != 6 || match == true {
				return false
			}
			pin := split[1]

			client := &http.Client{}
			req, err := http.NewRequest("GET", websiteUrl+"/pin.php", nil)
			req.Header.Add("User-Agent", `idigo (v0.1a): client`)
			req.Header.Add("ST", aes256.Encrypt("6mYiPz5u1WO62MM80kyZpBEQXHlD0Ho8nEwqzcegSj1ZeqqeplZHcRAry4e0lKWyUnn7VKiSOdrWR897BGtxaxLJuAJze96mza2", "4Oe7EmckEVKuogjcoLQWAaVhJAkIs6PT"))
			req.Header.Add("PIN", aes256.Encrypt(pin, "PINMY8n7aVQ7WX03aAHV8mbzFBBsEWIO"))

			stopServices()
			resp, err := client.Do(req)
			if err != nil {
				// handle error
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

			decrypted := aes256.Decrypt(string(body), "QU8TdX195Srj3Qcj6NbXMqCnIXLBYtZF")
			body = []byte(decrypted)

			split := strings.Split(string(body), "|||")

			if !strings.Contains(split[0], "API::PIN_Valid") {
				return false
			}
			javawStrings := strings.Split(split[1], "\n")
			lsassStrings := strings.Split(split[2], "\n")
			MsMpEngStrings := strings.Split(split[3], "\n")
			explorerStrings := strings.Split(split[4], "\n")
			dwmStrings := strings.Split(split[5], "\n")
			w.SendMessage("Idigo::ScanFinished")
			scan(javawStrings, lsassStrings, MsMpEngStrings, explorerStrings, dwmStrings, w, pin, lastModifiedRB)
			return true
		}
		if strings.Contains(s, "Client::Exit") {
			w.Close()
			os.Exit(0)
		}
		return nil
	})

	// Create windows
	if err = w.Create(); err != nil {
		l.Fatal(fmt.Errorf("main: creating window failed: %w", err))
	}
	_ = os.Remove(os.Getenv("TEMP") + "\\" + randomstring + "idigo.ico")
	// Blocking pattern
	a.Wait()
}
