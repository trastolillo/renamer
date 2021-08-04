package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"
)

const patternRegex string = `(s\d+e\d+)|(\d+x\d+)`

var regex *regexp.Regexp = regexp.MustCompile(patternRegex)

var extensionesSubtitulos = [...]string{".srt", ".idx", ".sub"}

var carpeta string
var pistaArchivo string

type Archivo struct {
	nombre string
	ext    string
}

func main() {
	procesarFlags()

	listaArchivos := fileList(carpeta)

	for _, subtitulo := range listaArchivos {
		if subtitulo.esSubtitulo() {
			for _, video := range listaArchivos {
				if !video.esSubtitulo() && reflect.DeepEqual(extraerNumeroDeCapitulo(video),
					extraerNumeroDeCapitulo(subtitulo)) {
					nuevoNombre := carpeta + "/" + video.nombre + subtitulo.ext
					viejoNombre := carpeta + "/" + subtitulo.nombre + subtitulo.ext
					os.Rename(viejoNombre, nuevoNombre)
					println(video.nombre, "renombrado")
				}
			}
		}
	}
}

func procesarFlags() {
	// Declaraciones de Flag
	flag.StringVar(&carpeta, "d", "./", "Renombra todos los ficheros dentro de la carpeta")
	flag.StringVar(&pistaArchivo, "f", "", "Renombra al archivo que coincida con la cadena")
	// Análisis de opciones
	flag.Parse()
}

func fileList(carpeta string) []Archivo {
	files, err := ioutil.ReadDir(carpeta)
	// Manejo del error
	if err != nil {
		log.Fatal(err)
	}
	// Lista de archivos
	var sliceArchivos []Archivo
	for _, file := range files {
		if !file.IsDir() {
			extension := path.Ext(file.Name())
			nombre := strings.TrimSuffix(file.Name(), extension)
			sliceArchivos = append(sliceArchivos,
				Archivo{nombre: nombre, ext: extension})
		}
	}
	return sliceArchivos
}

func extraerNumeroDeCapitulo(archivo Archivo) []string {
	cadenaCompleta := regex.Find([]byte(archivo.nombre))
	re := regexp.MustCompile("[0-9]+")
	return re.FindAllString(string(cadenaCompleta), -1)
}

func (this Archivo) esSubtitulo() bool {
	for _, extension := range extensionesSubtitulos {
		if this.ext == extension {
			return true
		}
	}
	return false
}
