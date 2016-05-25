package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
	"math/rand"
	"path/filepath"
	"github.com/kardianos/osext"
	"github.com/BurntSushi/toml"
	"os"
	"io"
)

const (
	//OFFICIAL_SITE = "http://pkg.entware.net/binaries/"
	TYPE_D = 1 // directory
	TYPE_F = 2 // file
)

var (
	Fetch_Delay = time.Millisecond * 400
	Config tomlConfig
)

type Msg struct  {
	Url string
	Type int
}

func fetch_detail(ch chan *Msg, rootUrl string, downloadFolder string) {
	for {
		data := <-ch

		file := strings.Replace(data.Url, rootUrl, downloadFolder, 1);

		if _, err := os.Stat(file); err == nil {
			log.Printf("File: %s exists, skip download, \n**WARNING**, We don't promise the file is downloaded complete but just check the file exists or not!*", file)
			continue
		}
		delay := get_delay_duration();
		if Config.Global.Debug { log.Printf("Delay %s to download url = %s", delay.String(), data.Url)}
		time.Sleep(delay)

		folder := file[0:strings.LastIndex(file,"/")]
		//log.Println("local folder = " + folder)

		if err := os.MkdirAll(folder, 0755); err != nil {
			log.Printf("make folder = %s failed! re-add to channel", folder)
			ch <- data
		} else {
			if resp, err := http.Get(data.Url); err != nil {
				log.Printf("download: %s error: %s, re-add", data.Url, err.Error())
				ch <- data
			} else {
				defer resp.Body.Close()
				if f, err := os.Create(file); err != nil {
					log.Printf("create file: %s failed, re-add", file)
					ch <- data
				} else {
					if length, err := io.Copy(f, resp.Body); err != nil {
						log.Printf("copy stream from url:%s failed, re-add", data.Url)
						ch <- data
					} else {
						if Config.Global.Debug { log.Printf("Download %d from url: %s", length, data.Url) }
					}
				}
			}
		}


	}
}

// get root folder
func fetch_list(ch chan *Msg, url string) {
	//time.Sleep(get_delay_duration())
	if Config.Global.Debug { log.Println("fetch url list, root = " + url) }

	if resp, err := http.Get(url); err != nil {
		log.Println(err.Error())
	} else {
		defer resp.Body.Close();
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			log.Println(err.Error())
		} else {
			content := string(body)
			//log.Println(content)

			regex := regexp.MustCompile("<a href=\"([^\\?\"]{2,})\">([^<]+)</a>");
			data := regex.FindAllStringSubmatch(content, -1)
			for i :=0; i < len(data); i++ {
				//log.Println(len(data[i]))
				if len(data[i]) == 3 {
					link := data[i][1]
					//log.Println(link)
					// skip un-used link
					if !strings.HasPrefix(link, "/") {
						//log.Println("link = " + link)
						if strings.HasSuffix(link, "/") {
							link = url + link
							if strings.Contains(link, "arm") ||
								strings.Contains(link, "others") ||
										strings.Contains(link, "mips") {

								fetch_list(ch, link)
							}
						} else {
							d := &Msg{
								Url: url + data[i][1],
							}
							//log.Println("file = " + d.Url)
							d.Type = TYPE_F
							ch <- d
						}
					}

				}
			}
		}
	}
}

func get_delay_duration() time.Duration {
	rand.Seed(time.Now().UnixNano())
	delay_duration := Fetch_Delay + time.Duration(time.Millisecond * time.Duration(rand.Int31n(500)))
	//log.Println("delay duration = " + delay_duration.String())
	return delay_duration
}

func main() {
	execPath, err := osext.ExecutableFolder()

	if err != nil {
		log.Fatal(err)
	}

	if _, err := toml.DecodeFile(execPath + string(filepath.Separator) + "config.toml", &Config); err != nil {
		log.Fatal(err)
	}

	ch := make(chan *Msg, 100000)
	go fetch_detail(ch, Config.Global.Root_url, Config.Global.Download_folder)
	fetch_list(ch, Config.Global.Root_url)
	for {
		log.Println("waiting...")
		time.Sleep(time.Second * 5)
	}
}
