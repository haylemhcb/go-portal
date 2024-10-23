package main

import (
    "net/http"
    "log"
    "os/exec"
    "encoding/json"
    "fmt"
    "os"
)

var State struct {
  isRun uint8
}

func deactivateMon() {

  fmt.Println("Previniendo modo monitor...")
  exec.Command("bash", "-c", "./unmon").Run()
  fmt.Println("OK")
}

func EscribirEnArchivo(nombreArchivo string, contenido string) error {
    // Abrir el archivo para escritura, creando uno nuevo si no existe
    // O sobrescribiendo el archivo existente
    archivo, err := os.Create(nombreArchivo)
    if err != nil {
        return err // Retorna un error si no se puede crear o abrir el archivo
    }
    defer archivo.Close() // Asegurarse de cerrar el archivo al final

    // Escribir el contenido en el archivo
    _, err = archivo.WriteString(contenido)
    if err != nil {
        return err // Retorna un error si no se puede escribir en el archivo
    }

    return nil // Retorna nil si todo fue exitoso
}


func LeerArchivo(ruta string) (string) {
    // Leer el contenido del archivo
    contenido, err := os.ReadFile(ruta)
    if err != nil {
        return ""
    }
    return string(contenido)
}


func readWopen(w http.ResponseWriter, r *http.Request) {
    exec.Command("bash", "-c", "./scanaps").Run()
}

func readWifisOpen(w http.ResponseWriter, r *http.Request) {
    exec.Command("bash", "-c", "./listaps").Run()

    lastw := LeerArchivo("./lastwifis")

    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprintln(w, lastw)

}


func readScap(w http.ResponseWriter, r *http.Request) {
    exec.Command("bash", "-c", "./readsignalcap").Run()

    scap := LeerArchivo("/tmp/scap")

    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprintln(w, scap)

}


func readStatus(w http.ResponseWriter, r *http.Request) {
    data := LeerArchivo("status")
    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprintln(w, data)
}

func readSignal(w http.ResponseWriter, r *http.Request) {
    exec.Command("bash", "-c", "./extractsign").Run();

    signal := LeerArchivo("/tmp/signal")

    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprintln(w, signal)

}


func stopRobot(w http.ResponseWriter, r *http.Request) {
    exec.Command("pkill", "airodump").Run()
    exec.Command("pkill", "dhclient").Run()
    exec.Command("pkill", "ping").Run()
    exec.Command("pkill", "portalrobot").Run()

    exec.Command("systemctl", "start", "NetworkManager").Run()
    exec.Command("rm", "-v", "./status").Run()
    exec.Command("rm", "-v", "/tmp/scap").Run()
    exec.Command("rm", "-v", "-r", "/tmp/robotcap").Run()

    os.Remove("/tmp/scap");
    os.Remove("./status");


    EscribirEnArchivo("status", "Detenido!!!")
    deactivateMon()

    State.isRun = 0;
}

func runRobot(w http.ResponseWriter, r *http.Request) {
    // Comando a ejecutar (cambia "ls" por "dir" si estás en Windows)
    //cmd := exec.Command("ls", "-l") // Para Windows, usa exec.Command("cmd", "/C", "dir")

    if State.isRun == 1 {
      return
    }

    State.isRun = 1;

    exec.Command("pkill", "airodump").Run()
    exec.Command("pkill", "dhclient").Run()
    exec.Command("pkill", "ping").Run()
    exec.Command("pkill", "portalrobot").Run()
    exec.Command("systemctl", "stop", "NetworkManager").Run()

    exec.Command("./portalrobot").Run();

}

func saveInterfaceHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }

    var data struct {
        Interface string `json:"interface"`
        Ssid string `json:"ssid"`
        Rate string `json:"rate"`
        Ap string `json:"ap"`
        Channel string `json:"channel"`
        Fakemac string `json:"fakemac"`
    }

    // Decodificar el JSON recibido
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
        return
    }

    file, err := os.OpenFile("robot.cfg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }

    defer file.Close()

    // Escribir los datos en el archivo
    if _, err := file.WriteString("interface=" + data.Interface + "\n"); err != nil {
        fmt.Println("Error al escribir en el archivo:", err)
        return
    }

    EscribirEnArchivo("interface_act", data.Interface)

    if _, err := file.WriteString("ssid=" + data.Ssid + "\n"); err != nil {
        fmt.Println("Error al escribir en el archivo:", err)
        return
    }


    if _, err := file.WriteString("rate=" + data.Rate + "\n"); err != nil {
        fmt.Println("Error al escribir en el archivo:", err)
        return
    }

    if _, err := file.WriteString("ap=" + data.Ap + "\n"); err != nil {
        fmt.Println("Error al escribir en el archivo:", err)
        return
    }

    if _, err := file.WriteString("channel=" + data.Channel + "\n"); err != nil {
        fmt.Println("Error al escribir en el archivo:", err)
        return
    }

    if _, err := file.WriteString("fakemac=" + data.Fakemac + "\n"); err != nil {
        fmt.Println("Error al escribir en el archivo:", err)
        return
    }

    if _, err := file.WriteString("gateway=" + "0.0.0.0" + "\n"); err != nil {
        fmt.Println("Error al escribir en el archivo:", err)
        return
    }

    // Guardar la interfaz en un archivo
    //ioutil.WriteFile("interfaz_seleccionada.txt", []byte("interface=" + data.Interface), 0644)
    //ioutil.WriteFile("interfaz_seleccionada.txt", []byte("ssid=" + data.Ssid), 0644)

    fmt.Fprintln(w, "Interfaz guardada exitosamente")
}

func interfacesHandler(w http.ResponseWriter, r *http.Request) {
    // Ejecutar el script Bash para obtener las interfaces
    cmd := exec.Command("bash", "-c", "./listintf")
    output, err := cmd.Output()
    if err != nil {
        http.Error(w, "Error al ejecutar listintf", http.StatusInternalServerError)
        return
    }

    // Devolver la salida como texto plano
    w.Header().Set("Content-Type", "text/plain")
    w.Write(output)
}

func main() {


    exec.Command("systemctl", "stop", "NetworkManager").Run()
    exec.Command("rm", "-v", "./status").Run()
    exec.Command("rm", "-v", "/tmp/scap").Run()
    exec.Command("rm", "-v", "-r", "/tmp/robotcap").Run()

    os.Remove("/tmp/scap");
    os.Remove("./status");
    deactivateMon()

    // Define la ruta de la carpeta pública
    fs := http.FileServer(http.Dir("./public"))

    // Maneja las solicitudes a la ruta raíz
    http.Handle("/", fs)
    http.HandleFunc("/interfaces", interfacesHandler)
    http.HandleFunc("/save-interface", saveInterfaceHandler)
    http.HandleFunc("/start", runRobot)
    http.HandleFunc("/stop", stopRobot)
    http.HandleFunc("/status", readStatus)
    http.HandleFunc("/signal", readSignal)
    http.HandleFunc("/lstwopen", readWopen)
    http.HandleFunc("/lstwifisopen", readWifisOpen)
    http.HandleFunc("/scap", readScap)

    // Inicia el servidor en el puerto 8080
    log.Println("Servidor escuchando en http://localhost:9742")
    err := http.ListenAndServe(":9742", nil)
    if err != nil {
        log.Fatal(err)
    }

    os.Remove("./status");
    defer deactivateMon()
    defer exec.Command("systemctl", "start", "NetworkManager").Run()

}
