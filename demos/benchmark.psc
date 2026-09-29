Algoritmo BenchmarkCompilador
    Definir i, j, n Como Entero
    Definir suma, acumulador Como Real

    n = 10000
    suma = 0

    i = 1
    Mientras i <= n Hacer
        acumulador = 0

        j = 1
        Mientras j <= 1000 Hacer
            acumulador = acumulador + (i * j) / (j + 1)

            Si j % 2 == 0 Entonces
                acumulador = acumulador + 1.5
            SiNo
                acumulador = acumulador - 0.5
            FinSi

            j = j + 1
        FinMientras

        suma = suma + acumulador
        i = i + 1
    FinMientras

    Escribir "Resultado: ", suma
FinAlgoritmo
