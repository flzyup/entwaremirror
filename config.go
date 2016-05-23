package main

type tomlConfig struct {
	Title  string
	Global global
}

type global struct {
	Debug bool
	Root_url string
	Download_folder string
}