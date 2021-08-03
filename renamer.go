package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

var carpeta string
var pistaArchivo string

type Archivo struct {
	nombre string
	ext    string
}

func main() {
	procesarFlags()

	listaArchivos := fileList(carpeta)

	for index, archivo := range listaArchivos {
		fmt.Println(index, archivo.nombre, archivo.ext)
	}

}

func procesarFlags() {
	// Declaraciones de Flag
	flag.StringVar(&carpeta, "d", "./", "Renombra todos los ficheros dentro de la carpeta")
	flag.StringVar(&pistaArchivo, "f", "", "Renombra al archivo que coincida con la cadena")
	// Análisis de opciones
	flag.Parse()

	fmt.Println(carpeta, pistaArchivo)

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
			sliceArchivos = append(sliceArchivos, Archivo{nombre: nombre, ext: extension})
		}
	}
	return sliceArchivos
}

func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}
