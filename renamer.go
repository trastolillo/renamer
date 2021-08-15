package main

import (
	"errors"
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
const appName = "Sub Renamer"

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
	// Validación
	if carpeta == "" {
		flag.Usage()
		os.Exit(2)
	}
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
	lista, fallo := listarArchivos(files, cadenaArchivo)
	if fallo != nil {
		log.Fatal(fallo)
	}
	return lista
}

func (this *Archivo) extraerNumeroDeCapitulo() {
	//? Extracción de temporada y episodio
	cadenaCompleta := regex.Find([]byte(this.nombre))
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
	fmt.Println("muestra", this.muestra)
}

func (this Archivo) compareTo(archivo Archivo) bool {
	var resultado bool = false
	if this.temporada != 0 {
		resultado = this.temporada == archivo.temporada &&
			this.episodio == archivo.episodio &&
			this.muestra == archivo.muestra &&
			this.nombre != archivo.nombre
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

func listarArchivos(files []fs.FileInfo, nombreDeArchivo string) ([]Archivo, error) {
	var sliceArchivos []Archivo
	condicion := true
	for _, file := range files {
		// fmt.Println("listarArchivos - variable file:", file.Name())
		fmt.Println("listarArchivos - variable nombreDeArchivo:", nombreDeArchivo)
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
	if len(sliceArchivos) == 0 {
		return nil, errors.New("No se encuentra el archivo")
	}
	return sliceArchivos, nil
}

func renombrar(listaArchivos *[]Archivo) {
	for _, subtitulo := range *listaArchivos {
		if subtitulo.esTipoArchivo(extensionesSubtitulos) {
			for _, video := range *listaArchivos {
				subtitulo.extraerNumeroDeCapitulo()
				video.extraerNumeroDeCapitulo()
				if video.esTipoArchivo(extensionesVideos) && subtitulo.compareTo(video) {
					nuevoNombre := carpeta + "/" + video.nombre + subtitulo.ext
					viejoNombre := carpeta + "/" + subtitulo.nombre + subtitulo.ext
					os.Rename(viejoNombre, nuevoNombre)
					println(viejoNombre, "renombrado a", nuevoNombre)
				} else {
					fmt.Println("Sin entrar al bucle de renombre")
				}
			}
		} else {
			fmt.Printf("%v%v no procesado\n", subtitulo.nombre, subtitulo.ext)
		}
	}
}

func usage() {
	msg := fmt.Sprintf(`usage: %s [OPTIONS]
	%s is a simple tool to rename subtitles files
	`, appName, appName)
	fmt.Println(msg)
	flag.PrintDefaults()
}
