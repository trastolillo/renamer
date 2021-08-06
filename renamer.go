package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

const patternRegex string = `([sS]\d+[eE]\d+)|(\d+[xX]\d+)`

var regex *regexp.Regexp = regexp.MustCompile(patternRegex)

var extensionesSubtitulos = []string{".srt", ".idx", ".sub"}
var extensionesVideos = []string{".avi", ".mkv", ".mpg", ".mp4", ".webm", ".wmv"}

var carpeta string
var cadenaArchivo string

type Archivo struct {
	nombre    string
	ext       string
	muestra   string
	temporada byte
	episodio  byte
}

func main() {
	procesarFlags()

	listaArchivos := fileList(carpeta)

	renombrar(&listaArchivos)

}

func procesarFlags() {
	// Declaraciones de Flag
	flag.StringVar(&carpeta, "d", "./", "Renombra todos los ficheros dentro de la carpeta")
	flag.StringVar(&cadenaArchivo, "f", "", "Renombra al archivo que coincida con la cadena")
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
	return listarArchivos(files, cadenaArchivo)
}

func (this *Archivo) extraerNumeroDeCapitulo() {
	//? Extracción de temporada y episodio
	cadenaCompleta := regex.Find([]byte(this.nombre))
	println("cadenacompleta", string(cadenaCompleta))
	re := regexp.MustCompile("[0-9]+")
	capString := re.FindAllString(string(cadenaCompleta), -1)
	if len(capString) > 0 {
		temporada, _ := strconv.ParseUint(capString[0], 10, 8)
		episodio, _ := strconv.ParseUint(capString[1], 10, 8)
		this.temporada = byte(temporada)
		this.episodio = byte(episodio)
	}
	//? Extracción de la muestra
	reMuestra := regexp.MustCompile("^[a-zA-Z]+")
	capMuestra := reMuestra.Find([]byte(this.nombre))
	this.muestra = string(capMuestra)
}

func (this Archivo) compareTo(archivo Archivo) bool {
	var resultado bool = false
	if this.temporada != 0 {
		resultado = this.temporada == archivo.temporada &&
			this.episodio == archivo.episodio &&
			strings.ToLower(this.muestra) == strings.ToLower(this.muestra) &&
			strings.ToLower(this.nombre) != strings.ToLower(archivo.nombre)
	}
	return resultado
}

func (this Archivo) esTipoArchivo(listaExtensiones []string) bool {
	for _, extension := range listaExtensiones {
		if this.ext == extension {
			return true
		}
	}
	return false
}

func listarArchivos(files []fs.FileInfo, nombreDeArchivo string) []Archivo {
	var sliceArchivos []Archivo
	condicion := true
	for _, file := range files {
		if nombreDeArchivo != "" {
			condicion = strings.Contains(strings.ToLower(file.Name()), strings.ToLower(nombreDeArchivo))
		}
		if !file.IsDir() && condicion {
			extension := path.Ext(file.Name())
			nombre := strings.TrimSuffix(file.Name(), extension)
			sliceArchivos = append(sliceArchivos,
				Archivo{nombre: nombre, ext: extension})
		}
	}
	return sliceArchivos
}

func renombrar(listaArchivos *[]Archivo) {
	var contador byte = 0
	for _, subtitulo := range *listaArchivos {
		if subtitulo.esTipoArchivo(extensionesSubtitulos) {
			for _, video := range *listaArchivos {
				subtitulo.extraerNumeroDeCapitulo()
				video.extraerNumeroDeCapitulo()
				// fmt.Println("video", video)
				// fmt.Println("sub", subtitulo)
				if video.esTipoArchivo(extensionesVideos) && subtitulo.compareTo(video) {
					nuevoNombre := carpeta + "/" + video.nombre + subtitulo.ext
					viejoNombre := carpeta + "/" + subtitulo.nombre + subtitulo.ext
					os.Rename(viejoNombre, nuevoNombre)
					println(viejoNombre, "renombrado a", nuevoNombre)
					contador++
				} else {
					fmt.Println("Sin entrar al bucle de renombre")
				}
			}
		}
	}
	if contador == 0 {
		fmt.Println("Ningún archivo renombrado")
	} else {
		fmt.Println("Archivos renombrados:", contador)
	}
}
