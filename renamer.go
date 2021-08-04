package main

import (
	"errors"
	"flag"
	"fmt"
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
		// subtitulo, err := seleccionarArchivoSubtitulo(&listaArchivos, subtitulo)
		// if err != nil {
		// 	fmt.Printf("%v/%v%v\n", carpeta, subtitulo.nombre, subtitulo.ext)
		// 	remove(listaArchivos, index)
		// }
		if subtitulo.esSubtitulo() {
			for _, video := range listaArchivos {
				if !video.esSubtitulo() && reflect.DeepEqual(extraerNumeroDeCapitulo(video),
					extraerNumeroDeCapitulo(subtitulo)) {
					nuevoNombre := carpeta + "/" + video.nombre + subtitulo.ext
					viejoNombre := carpeta + "/" + subtitulo.nombre + subtitulo.ext
					println(nuevoNombre)
					println(viejoNombre)
					os.Rename(viejoNombre, nuevoNombre)
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

func seleccionarArchivoSubtitulo(videos *[]Archivo, subtitulo Archivo) (Archivo, error) {
	err := errors.New("No hay coincidencia entre el archivo de vídeo y el de subtítulos")
	vacio := Archivo{}
	for index, video := range *videos {
		fmt.Println(video.nombre, subtitulo, subtitulo.esSubtitulo() && !video.esSubtitulo() && reflect.DeepEqual(extraerNumeroDeCapitulo(video),
			extraerNumeroDeCapitulo(subtitulo)))
		if subtitulo.esSubtitulo() && !video.esSubtitulo() && reflect.DeepEqual(extraerNumeroDeCapitulo(video),
			extraerNumeroDeCapitulo(subtitulo)) {
			entrega := subtitulo
			fmt.Println("dentro: ", subtitulo, "long", len(*videos))
			remove(*videos, index)
			return entrega, nil
		}
	}
	return vacio, err
}

func (this Archivo) esSubtitulo() bool {
	for _, extension := range extensionesSubtitulos {
		if this.ext == extension {
			return true
		}
	}
	return false
}

func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func remove(a []Archivo, i int) []Archivo {
	copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index.
	// a[len(a)-1] = 0      // Erase last element (write zero value).
	// a = a[:len(a)-1]     // Truncate slice.
	b := a[:len(a)-1]
	return b
}
