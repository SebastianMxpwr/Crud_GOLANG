package main

import (
	// "fmt"
	"database/sql"
	"log"
	"net/http"
	"text/template"

	//para usar un driver usamos el guion bajo
	_ "github.com/go-sql-driver/mysql"
)

//el segundo parentesis es que retorna
func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := "12345SMM"
	Nombre := "crud_go"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

//busca las plantillas o el folder con el nombre
var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", start)
	http.HandleFunc("/crear", crear)
	http.HandleFunc("/insertar", insert)
	http.HandleFunc("/borrar", borrar)
	http.HandleFunc("/editar", editar)
	http.HandleFunc("/actulizar", actualizar)
	log.Println("Servidor Corriendo")
	http.ListenAndServe(":8080", nil)
}

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func borrar(w http.ResponseWriter, r *http.Request) {
	//obtener datos de los paramatros
	idEmpleado := r.URL.Query().Get("id")
	log.Println(idEmpleado)

	conexionEstablecida := conexionDB()
	borrarRegistro, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	borrarRegistro.Exec(idEmpleado)
	http.Redirect(w, r, "/", 301)
}

func editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	log.Println(idEmpleado)

	conexionEstablecida := conexionDB()
	registro, err := conexionEstablecida.Query("SELECT * FROM empleados WHERE id =?", idEmpleado)
	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}

	for registro.Next() {
		var id int
		var nombre, correo string
		err = registro.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

	}

	log.Println(empleado)
	plantillas.ExecuteTemplate(w, "editar", empleado)

}

func actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")
		conexionEstablecida := conexionDB()
		modificarRegistros, err := conexionEstablecida.Prepare("UPDATE empleados SET nombre=?, correo=? WHERE id =?")
		if err != nil {
			panic(err.Error())
		}
		modificarRegistros.Exec(nombre, correo, id)

		http.Redirect(w, r, "/", 301)

	}
}

func start(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionDB()
	registros, err := conexionEstablecida.Query("SELECT * FROM empleados")
	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}
	arrEmpleados := []Empleado{}

	for registros.Next() {
		var id int
		var nombre, correo string
		err = registros.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arrEmpleados = append(arrEmpleados, empleado)
	}

	// fmt.Println(arrEmpleados)
	//hace alusion a la plantilla el unico string
	plantillas.ExecuteTemplate(w, "inicio", arrEmpleados)

}

func crear(w http.ResponseWriter, r *http.Request) {

	plantillas.ExecuteTemplate(w, "crear", nil)

}

func insert(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")
		conexionEstablecida := conexionDB()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre,correo) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, correo)

		http.Redirect(w, r, "/", 301)

	}
}
