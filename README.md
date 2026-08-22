<div align="center">

# PASC

**Compilador experimental de pseudocódigo PSeInt a ejecutables nativos mediante Go.**

[![Go](https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Estado](https://img.shields.io/badge/estado-en%20desarrollo-f59e0b)](#estado-del-proyecto)
[![Plataforma](https://img.shields.io/badge/plataforma-multiplataforma-3c3c3c)](#requisitos)
[![Código fuente](https://img.shields.io/badge/GitHub-sea2horses%2Fpseint--compiler-181717?logo=github)](https://github.com/sea2horses/pseint-compiler)

</div>

PASC toma un archivo `.psc`, comprueba su sintaxis y sus tipos, lo traduce a Go y genera un binario nativo. Está pensado como una base moderna para ejecutar pseudocódigo de estilo PSeInt con mejores diagnósticos y aprovechando el compilador de Go como backend.

> [!IMPORTANT]
> El proyecto está en una etapa temprana (`Indev`). Ya puede ejecutar programas útiles, pero todavía no implementa todo el lenguaje de PSeInt ni debe considerarse compatible al 100 %.

## Lo mejor del proyecto

- **Compilación real:** el pseudocódigo se transforma en Go y después en un ejecutable nativo.
- **Análisis por etapas:** incluye lexer, parser, árbol sintáctico tipado, análisis semántico y generación de código.
- **Errores fáciles de ubicar:** muestra el archivo, la línea, la columna, el fragmento afectado y un resumen de errores y advertencias.
- **Tipado y conversiones:** valida operaciones y asignaciones antes de generar el programa, con soporte para conversiones numéricas seguras.
- **Sintaxis familiar:** reconoce palabras clave sin distinguir mayúsculas y minúsculas.
- **Entrada y salida interactiva:** permite leer y escribir varios valores, además de limpiar la terminal.
- **Modo de depuración:** expone tokens, árboles y código generado para facilitar el desarrollo del compilador.

## Funcionalidades disponibles

Actualmente se soportan:

- Programas delimitados por `Algoritmo` y `FinAlgoritmo`.
- Tipos `Entero`, `Real`, `Cadena` y `Logico`.
- Declaraciones explícitas y asignaciones.
- Literales numéricos, cadenas y valores lógicos.
- `Escribir`, `Leer` y `Borrar Pantalla`.
- Condicionales `Si`, `SiNo` y `FinSi`.
- Ciclos `Mientras ... Hacer` y `FinMientras`.
- Selección con `Segun`, casos, `De Otro Modo` y `FinSegun`.
- Operadores aritméticos, relacionales y lógicos, incluyendo módulo.
- Paréntesis y precedencia de operadores.

## Requisitos

- [Go 1.25 o posterior](https://go.dev/dl/).
- Git, únicamente si vas a clonar el repositorio.

El compilador invoca a Go internamente para construir el programa generado, por lo que `go` debe estar disponible en el `PATH` incluso después de compilar esta herramienta. (Mejora que se debe hacer)

## Instalación

Clona el repositorio y compila la herramienta:

```bash
git clone https://github.com/sea2horses/pseint-compiler.git
cd pseint-compiler
go build -o pseintc .
```

También puedes usarla directamente desde el código fuente con `go run .`.

## Uso

Ejecuta el compilador indicando un archivo `.psc` con `-file`:

```bash
./pseintc -file demos/hello_world.psc
```

O, sin construir primero el compilador:

```bash
go run . -file demos/hello_world.psc
```

Por el momento, cada ejecución realiza todo el proceso:

1. Lee y valida el pseudocódigo.
2. Genera `build/out.go`.
3. Compila el ejecutable `build/out` (`build/out.exe` en Windows).
4. Ejecuta el binario automáticamente.

### Ejemplo

Crea un archivo llamado `mayoria_edad.psc`:

```text
Algoritmo MayoriaDeEdad
    Definir edad Como Entero

    Escribir "¿Cuántos años tienes?"
    Leer edad

    Si edad >= 18 Entonces
        Escribir "Eres mayor de edad"
    SiNo
        Escribir "Eres menor de edad"
    FinSi
FinAlgoritmo
```

Compílalo y ejecútalo:

```bash
./pseintc -file mayoria_edad.psc
```

### Información de depuración

Usa `-debug` para ver las etapas internas, como los tokens, el AST, el árbol tipado y el código Go generado:

```bash
./pseintc -debug -file demos/operations.psc
```

## Estado del proyecto

La base del compilador ya funciona de extremo a extremo, pero aún faltan piezas importantes para alcanzar una experiencia completa y una compatibilidad amplia con PSeInt:

- Añadir ciclos `Para` y `Repetir ... Hasta Que`.
- Implementar subprocesos, funciones y retorno de valores.
- Incorporar arreglos/dimensiones y acceso por índices.
- Agregar caracteres y más variantes de la sintaxis de PSeInt.
- Completar el operador de potencia, que el lexer reconoce pero el análisis semántico todavía no admite.
- Permitir compilar sin ejecutar, elegir el nombre y la ubicación del binario, y mostrar una opción de versión o ayuda más completa.
- Manejar de forma explícita errores y entradas largas durante `Leer`.
- Ampliar y estabilizar las reglas de alcance, asignación implícita y conversiones numéricas.
- Crear una suite automatizada de pruebas unitarias, de integración y de programas `.psc`.
- Publicar versiones etiquetadas, binarios precompilados y una licencia para el proyecto.

## Probar los ejemplos

El directorio `demos` contiene programas para probar entrada/salida, operaciones, condicionales, ciclos y `Segun`. Por ejemplo:

```bash
go run . -file demos/operations.psc
go run . -file demos/ages.psc
go run . -file demos/switch.psc
```

Para comprobar que todos los paquetes compilan:

```bash
go test ./...
```

Actualmente no hay casos de prueba automatizados; este comando funciona como comprobación de compilación hasta que se incorpore la suite de pruebas.

## Contribuciones

Las contribuciones son bienvenidas, especialmente para ampliar la sintaxis compatible, mejorar los diagnósticos y añadir pruebas. Antes de hacer un cambio grande, conviene abrir un issue para acordar el alcance y evitar implementaciones duplicadas.