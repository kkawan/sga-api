package handlers

import (
	"errors"
	"strconv"
	"strings"
)

// dias que o sistema aceita
var diasDaSemana = []string{"segunda", "terca", "quarta", "quinta", "sexta", "sabado", "domingo"}

// arrumaDia deixa o dia sempre no mesmo padrão pra poder comparar depois
func arrumaDia(dia string) string {
	return strings.ToLower(strings.TrimSpace(dia))
}

// diaValido confere se o dia existe na lista
func diaValido(dia string) bool {
	for _, d := range diasDaSemana {
		if d == dia {
			return true
		}
	}
	return false
}

// horaParaMinutos transforma "14:30" em 870 minutos. Fica bem mais fácil
// comparar dois horários usando número do que usando texto.
func horaParaMinutos(hora string) (int, error) {
	partes := strings.Split(strings.TrimSpace(hora), ":")
	if len(partes) != 2 {
		return 0, errors.New("o horário precisa estar no formato HH:MM")
	}

	h, err := strconv.Atoi(partes[0])
	if err != nil || h < 0 || h > 23 {
		return 0, errors.New("hora inválida")
	}

	m, err := strconv.Atoi(partes[1])
	if err != nil || m < 0 || m > 59 {
		return 0, errors.New("minuto inválido")
	}

	return (h * 60) + m, nil
}

// temSobreposicao devolve true quando os dois horários se cruzam.
// A regra do enunciado é: inicio novo < fim existente E fim novo > inicio existente.
func temSobreposicao(inicioNovo, fimNovo, inicioAntigo, fimAntigo int) bool {
	return inicioNovo < fimAntigo && fimNovo > inicioAntigo
}
