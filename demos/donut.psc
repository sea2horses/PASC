Algoritmo Dona3D
	Definir A, B, i, j Como Real
	Definir c, d, e, f, g, h Como Real
	Definir DD, l, m, n, t Como Real
	Definir x, yy, oo, NN Como Entero

	Dimensionar bb[1760]
	Dimensionar z[1760]

	A = 0
	B = 0

	Mientras Verdadero Hacer

		i = 0
		Mientras i <= 1759 Hacer
			bb[i] = " "
			z[i] = 0
			i = i + 1
		FinMientras

		j = 0

		Mientras j < 6.28 Hacer
			i = 0

			Mientras i < 6.28 Hacer

				c = Sen(i)
				d = Cos(j)
				e = Sen(A)
				f = Sen(j)
				g = Cos(A)
				h = d + 2

				DD = 1 / (c * h * e + f * g + 5)

				l = Cos(i)
				m = Cos(B)
				n = Sen(B)

				t = c * h * g - f * e

				x = Trunc(40 + 30 * DD * (l * h * m - t * n))
				yy = Trunc(12 + 15 * DD * (l * h * n + t * m))

				oo = x + 80 * yy

				NN = Trunc(8 * ((f * e - c * d * g) * m - c * d * e - f * g - l * d * n))

				Si yy >= 0 && yy < 22 && x >= 0 && x < 80 Entonces
					Si DD > z[oo] Entonces
						z[oo] = DD

						Si NN < 0 Entonces
							NN = 0
						FinSi

						Si NN > 11 Entonces
							NN = 11
						FinSi

						Segun NN Hacer
							0:
								bb[oo] = "."
							1:
								bb[oo] = ","
							2:
								bb[oo] = "-"
							3:
								bb[oo] = "~"
							4:
								bb[oo] = ":"
							5:
								bb[oo] = ";"
							6:
								bb[oo] = "="
							7:
								bb[oo] = "!"
							8:
								bb[oo] = "*"
							9:
								bb[oo] = "#"
							10:
								bb[oo] = "$"
							11:
								bb[oo] = "@"
						FinSegun
					FinSi
				FinSi

				i = i + 0.02
			FinMientras

			j = j + 0.07
		FinMientras

		Borrar Pantalla

		yy = 0
		Mientras yy <= 21 Hacer

			x = 0
			Mientras x <= 79 Hacer
				Escribir Sin Saltar bb[x + 80 * yy]
				x = x + 1
			FinMientras

			Escribir ""
			yy = yy + 1
		FinMientras

		A = A + 0.04
		B = B + 0.02

	FinMientras

FinAlgoritmo
