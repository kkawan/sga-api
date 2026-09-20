package handlers

import (
	"errors"
	"strconv"
	"strings"
)

// dias que o sistema aceita. Coloquei sem acento pra facilitar a comparação,
// na hora de salvar eu sempre converto pra minúsculo.
var diasDaSemana = []string{"segunda", "terca", "quarta", "quinta", "sexta", "sabado", "domingo"}

// arrumaDia deixa o dia minúsculo e tira o acento do "terça" e do "sábado",
// senão "Terça" e "terca" iam ser tratados como dias diferentes.
func arrumaDia(dia string) string {
	d := strings.ToLower(strings.TrimSpace(dia))
	d = strings.ReplaceAll(d, "ç", "c")
	d = strings.ReplaceAll(d, "á", "a")
	d = strings.ReplaceAll(d, "-feira", "")
	return strings.TrimSpace(d)
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
