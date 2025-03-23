package main

import (
    "bufio"
    "encoding/base64"
    "fmt"
    "os"
    "os/exec"
    "regexp"
    "strconv"
    "strings"
    "sync"
    "time"
)

// Global vars
var (
    maxThreads  int = 1
    wg          sync.WaitGroup
    verbose     bool = false
    rateLimit   int = 5
    rateBoolean bool = true
    sem         = make(chan int, maxThreads)
    lastPart, previousParts string
    xV, xH, xUA, xX, xD, xS, xM, xE, xB bool = false, false, false, false, false, false, false, false, false
)

func banner() {
    fmt.Println("\033[1;36m") 
    fmt.Println("  ____  _ _   ______      _")
    fmt.Println(" |  _ \\(_) | |  _ \\ \\    / /")
    fmt.Println(" | |_) |_| |_| |_) \\ \\/\\/ / ")
    fmt.Println(" |  __/| | __|  __/ \\  / /  ")
    fmt.Println(" | |   | | |_| |     / /\\ \\ ")
    fmt.Println(" |_|   |_|\\__|_|    /_/  \\_\\")
    fmt.Println("\033[1;35m") 
    fmt.Println(" PiteBypass - Bypass 4xx like a pro!")
    fmt.Println("\033[1;33m") 
    fmt.Println(" by: \033[1;31m@pitecozz\033[0m")
    fmt.Println("\033[0m") 
}

func main() {
    banner()

    // Verifica se o primeiro argumento é -h ou --help
    if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
        showHelp()
        os.Exit(0)
    }

    // Verifica se o número de argumentos é menor que 2 (sem URL ou arquivo)
    if len(os.Args) < 2 {
        showHelp()
        os.Exit(0)
    }

    // O último argumento é a URL ou arquivo
    target := os.Args[len(os.Args)-1]
    options := os.Args[1 : len(os.Args)-1]

    processOptions(options)

    if isValidURL(target) {
        byp4xx(options, target)
    } else {
        processFile(target, options)
    }
}

func showHelp() {
    fmt.Println("byp4xx <cURL or byp4xx options> <URL or file>")
    fmt.Println("Some cURL options you may use as example:")
    fmt.Println("  -L follow redirections (30X responses)")
    fmt.Println("  -x <ip>:<port> to set a proxy")
    fmt.Println("  -m <seconds> to set a timeout")
    fmt.Println("  -H for new headers. Escape double quotes.")
    fmt.Println("  -d for data in the POST requests body")
    fmt.Println("  -...")
    fmt.Println("Built-in options:")
    fmt.Println("  --all Verbose mode")
    fmt.Println("  -t or --thread Set the maximum threads")
    fmt.Println("  --rate Set the maximum reqs/sec. Only one thread enforced, for low rate limits.")
    fmt.Println("  -xV Exclude verb tampering")
    fmt.Println("  -xH Exclude headers")
    fmt.Println("  -xUA Exclude User-Agents")
    fmt.Println("  -xX Exclude extensions")
    fmt.Println("  -xD Exclude default creds")
    fmt.Println("  -xS Exclude CaSe SeNsiTiVe")
    fmt.Println("  -xM Exclude middle paths")
    fmt.Println("  -xE Exclude end paths")
    fmt.Println("  -xB Exclude #bugbountytips")
}

func processOptions(options []string) {
    for i := 0; i < len(options); i++ {
        switch options[i] {
        case "--rate":
            if i+1 < len(options) {
                rateLimit, _ = strconv.Atoi(options[i+1])
                maxThreads = 1
                sem = make(chan int, maxThreads)
                options = append(options[:i], options[i+2:]...)
            }
        case "-t", "--thread":
            if i+1 < len(options) {
                maxThreads, _ = strconv.Atoi(options[i+1])
                sem = make(chan int, maxThreads)
                options = append(options[:i], options[i+2:]...)
                rateBoolean = false
            }
        case "--all":
            verbose = true
            options = append(options[:i], options[i+1:]...)
        case "-xV":
            xV = true
            options = append(options[:i], options[i+1:]...)
        case "-xH":
            xH = true
            options = append(options[:i], options[i+1:]...)
        case "-xUA":
            xUA = true
            options = append(options[:i], options[i+1:]...)
        case "-xX":
            xX = true
            options = append(options[:i], options[i+1:]...)
        case "-xD":
            xD = true
            options = append(options[:i], options[i+1:]...)
        case "-xS":
            xS = true
            options = append(options[:i], options[i+1:]...)
        case "-xM":
            xM = true
            options = append(options[:i], options[i+1:]...)
        case "-xE":
            xE = true
            options = append(options[:i], options[i+1:]...)
        case "-xB":
            xB = true
            options = append(options[:i], options[i+1:]...)
        }
    }
}

func isValidURL(url string) bool {
    match, _ := regexp.MatchString("^https?://", url)
    return match
}

func processFile(filePath string, options []string) {
    file, err := os.Open(filePath)
    if err != nil {
        fmt.Println("Error reading file:", err)
        os.Exit(1)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if isValidURL(line) {
            byp4xx(options, line)
        } else {
            fmt.Println("Invalid URL in file:", line)
        }
    }
}

func byp4xx(options []string, url string) {
    parts := strings.Split(strings.TrimRight(url, "/"), "/")
    lastPart = parts[len(parts)-1]
    previousParts = strings.Join(parts[:len(parts)-1], "/")

    fmt.Println("\033[31m===== " + url + " =====\033[0m")
    if !xV {
        verbTampering(options, url)
    }
    if !xH {
        headers(options, url)
    }
    if !xUA {
        userAgent(options, url)
    }
    if !xX {
        extensions(options, url)
    }
    if !xD {
        defaultCreds(options, url)
    }
    if !xS {
        caseSensitive(options, url)
    }
    if !xM {
        midPaths(options, url)
    }
    if !xE {
        endPaths(options, url)
    }
    if !xB {
        bugBounty(options, url)
    }
}

func curlCodeResponse(message string, options []string, url string) {
    codeOptions := []string{"-k", "-s", "-o", "/dev/null", "-w", "\"%{http_code}\""}
    payload := append(options, url)
    payload = append(codeOptions, payload...)
    curlCommand := exec.Command("curl", payload...)

    output, _ := curlCommand.CombinedOutput()
    outputStr := strings.ReplaceAll(string(output), "\"", "")
    code, _ := strconv.Atoi(outputStr)

    if code >= 200 && code < 300 {
        outputStr = "\033[32m" + outputStr + "\033[0m"
    } else if code >= 300 && code < 400 {
        outputStr = "\033[33m" + outputStr + "\033[0m"
    } else if verbose {
        outputStr = "\033[31m" + outputStr + "\033[0m"
    }

    fmt.Println(message, outputStr)
    if rateBoolean {
        rateLimitMod := 1.0 / float64(rateLimit) * 1000.0
        time.Sleep(time.Duration(rateLimitMod) * time.Millisecond)
    }
    wg.Done()
}

func verbTampering(options []string, url string) {
    fmt.Println("\033[32m==VERB TAMPERING==\033[0m")
    file, _ := os.Open("/usr/local/share/byp4xx/templates/verbs.txt")
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        sem <- 1
        wg.Add(1)
        go func(line string) {
            optionsMod := append(options, "-X", line)
            curlCodeResponse(line+":", optionsMod, url)
            <-sem
        }(line)
    }
    wg.Wait()
}

func headers(options []string, url string) {
    fmt.Println("\033[32m==HEADERS==\033[0m")
    file, err := os.Open("/usr/local/share/byp4xx/templates/headers.txt")
    if err != nil {
        fmt.Println("Error opening headers file:", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        header := scanner.Text()
        file2, err := os.Open("/usr/local/share/byp4xx/templates/ip.txt")
        if err != nil {
            fmt.Println("Error opening IP file:", err)
            return
        }
        defer file2.Close()

        scanner2 := bufio.NewScanner(file2)
        for scanner2.Scan() {
            ip := scanner2.Text()
            sem <- 1
            wg.Add(1)
            go func(header, ip string) {
                headerLine := header + ip
                optionsMod := append(options, "-H", headerLine)
                curlCodeResponse(headerLine+":", optionsMod, url)
                <-sem
            }(header, ip)
        }
    }
    wg.Wait()
}

func userAgent(options []string, url string) {
    fmt.Println("\033[32m==USER AGENTS==\033[0m")
    file, err := os.Open("/usr/local/share/byp4xx/templates/UserAgents.txt")
    if err != nil {
        fmt.Println("Error opening User-Agents file:", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        userAgent := scanner.Text()
        sem <- 1
        wg.Add(1)
        go func(userAgent string) {
            header := "User-Agent: " + userAgent
            optionsMod := append(options, "-H", header)
            curlCodeResponse(header+":", optionsMod, url)
            <-sem
        }(userAgent)
    }
    wg.Wait()
}

func extensions(options []string, url string) {
    fmt.Println("\033[32m==EXTENSIONS==\033[0m")
    file, err := os.Open("/usr/local/share/byp4xx/templates/extensions.txt")
    if err != nil {
        fmt.Println("Error opening extensions file:", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        extension := scanner.Text()
        sem <- 1
        wg.Add(1)
        go func(extension string) {
            urlMod := url + extension
            curlCodeResponse(extension+":", options, urlMod)
            <-sem
        }(extension)
    }
    wg.Wait()
}

func defaultCreds(options []string, url string) {
    fmt.Println("\033[32m==DEFAULT CREDS==\033[0m")
    file, err := os.Open("/usr/local/share/byp4xx/templates/defaultcreds.txt")
    if err != nil {
        fmt.Println("Error opening default credentials file:", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        creds := scanner.Text()
        sem <- 1
        wg.Add(1)
        go func(creds string) {
            sEnc := base64.StdEncoding.EncodeToString([]byte(creds))
            header := "Authorization: Basic " + sEnc
            optionsMod := append(options, "-H", header)
            curlCodeResponse(creds+":", optionsMod, url)
            <-sem
        }(creds)
    }
    wg.Wait()
}

func caseSensitive(options []string, url string) {
    fmt.Println("\033[32m==CASE SENSITIVE==\033[0m")
    for i := range lastPart {
        modifiedURI := ""
        for j, r := range lastPart {
            if j == i {
                if r >= 'A' && r <= 'Z' {
                    modifiedURI += string(r + ('a' - 'A'))
                } else if r >= 'a' && r <= 'z' {
                    modifiedURI += string(r - ('a' - 'A'))
                }
            } else {
                modifiedURI += string(r)
            }
        }
        sem <- 1
        wg.Add(1)
        go func(modifiedURI string) {
            urlMod := previousParts + "/" + modifiedURI
            curlCodeResponse(modifiedURI+":", options, urlMod)
            <-sem
        }(modifiedURI)
    }
    wg.Wait()
}

func midPaths(options []string, url string) {
    fmt.Println("\033[32m==MID PATHS==\033[0m")
    file, err := os.Open("/usr/local/share/byp4xx/templates/midpaths.txt")
    if err != nil {
        fmt.Println("Error opening midpaths file:", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        path := scanner.Text()
        sem <- 1
        wg.Add(1)
        go func(path string) {
            optionsMod := append(options, "--path-as-is")
            urlMod := previousParts + "/" + path + "/" + lastPart
            curlCodeResponse(path+":", optionsMod, urlMod)
            <-sem
        }(path)
    }
    wg.Wait()
}

func endPaths(options []string, url string) {
    fmt.Println("\033[32m==END PATHS==\033[0m")
    file, err := os.Open("/usr/local/share/byp4xx/templates/endpaths.txt")
    if err != nil {
        fmt.Println("Error opening endpaths file:", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        path := scanner.Text()
        sem <- 1
        wg.Add(1)
        go func(path string) {
            optionsMod := append(options, "--path-as-is")
            urlMod := previousParts + "/" + lastPart + path
            curlCodeResponse(path+":", optionsMod, urlMod)
            <-sem
        }(path)
    }
    wg.Wait()
}

func bugBounty(options []string, url string) {
    fmt.Println("\033[32m==BUG BOUNTY TIPS==\033[0m")
    tests := []struct {
        path    string
        options []string
    }{
        {"/%2e/" + lastPart, options},
        {"/%ef%bc%8f" + lastPart, append(options, "--path-as-is")},
        {"/" + lastPart + "?", options},
        {"/" + lastPart + "??", options},
        {"/" + lastPart + "//", options},
        {"/" + lastPart + "/", options},
        {"/./" + lastPart + "/./", append(options, "--path-as-is")},
        {"/" + lastPart + "/.randomstring", options},
        {"/" + lastPart + "..;/", append(options, "--path-as-is")},
        {"/" + lastPart + "..;", append(options, "--path-as-is")},
        {"/.;/" + lastPart, append(options, "--path-as-is")},
        {"/.;/" + lastPart + "/.;/", append(options, "--path-as-is")},
        {"/;foo=bar/" + lastPart, append(options, "--path-as-is")},
    }

    for _, test := range tests {
        sem <- 1
        wg.Add(1)
        go func(test struct {
            path    string
            options []string
        }) {
            urlMod := previousParts + test.path
            curlCodeResponse(test.path+":", test.options, urlMod)
            <-sem
        }(test)
    }
    wg.Wait()
}

