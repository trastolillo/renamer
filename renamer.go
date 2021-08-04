package main

import (
	"flag"
	"io/fs"
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

var extensionesSubtitulos = []string{".srt", ".idx", ".sub"}
var extensionesVideos = []string{".avi", ".mkv", ".mpg", ".mp4", ".webm", ".wmv"}

var carpeta string
var cadenaArchivo string

type Archivo struct {
	nombre string
	ext    string
}

func main() {
	procesarFlags()

	listaArchivos := fileList(carpeta)

	for _, subtitulo := range listaArchivos {
		if subtitulo.esTipoArchivo(extensionesSubtitulos) {
			for _, video := range listaArchivos {
				if video.esTipoArchivo(extensionesVideos) && reflect.DeepEqual(extraerNumeroDeCapitulo(video),
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

func extraerNumeroDeCapitulo(archivo Archivo) []string {
	cadenaCompleta := regex.Find([]byte(archivo.nombre))
	re := regexp.MustCompile("[0-9]+")
	return re.FindAllString(string(cadenaCompleta), -1)
}

func (this Archivo) esTipoArchivo(listaExtensiones []string) bool {
	for _, extension := range listaExtensiones {
		if this.ext == extension {
			return true
		}
	}
	return false
}

func filtrarArchivosPorNombre(listaArchivos []Archivo) []Archivo {
	// retorno := []Archivo{}
	return nil

}

func listarArchivos(files []fs.FileInfo, nombreDeArchivo string) []Archivo {
	var sliceArchivos []Archivo
	condicion := true
	for _, file := range files {
		if nombreDeArchivo != "" {
			condicion = strings.Contains(file.Name(), nombreDeArchivo)
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
