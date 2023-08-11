package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"text/template"
)

type Question struct {
	ID       int    `json:"id"`
	Question string `json:"question"`

	UserAnswer string `json:"user_answer"`
}

func main() {
	http.HandleFunc("/", questionario)
	http.HandleFunc("/executar", executarArquivoShell)
	http.ListenAndServe(":8081", nil)
}

func questionario(w http.ResponseWriter, r *http.Request) {

	templ := template.Must(template.ParseFiles("templates/index.html"))

	if r.Method == http.MethodPost {
		questions := []Question{
			{ID: 1, Question: "Qual o nome do subdomínio?", UserAnswer: r.FormValue("answer1")},
			{ID: 2, Question: "Qual o TTL em segundos?", UserAnswer: r.FormValue("answer2")},
			{ID: 3, Question: "Qual a nome da zona de cadastro?", UserAnswer: r.FormValue("answer3")},
			{ID: 4, Question: "Qual o Grupo de Recurso?", UserAnswer: r.FormValue("answer4")},
		}

		// Criando um arquivo para escrever os dados JSON
		file, err := os.Create("resultado.json")
		if err != nil {
			fmt.Println("Erro ao criar o arquivo:", err)
			return
		}
		defer file.Close()

		// Convertendo o questionário em formato JSON
		resultado, err := json.MarshalIndent(questions, "", "    ")
		if err != nil {
			fmt.Println("Erro ao converter para JSON:", err)
			return
		}

		// Escrevendo o JSON no arquivo
		_, err = file.Write(resultado)
		if err != nil {
			fmt.Println("Erro ao escrever no arquivo:", err)
			return
		}

		fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			
		</head>
		<body>
			<form action="executar">
				<button>EXECUTAR</button>
			</form>
		</body>
		<script>
			document.getElementById('meuFormulario').addEventListener('submit', function() {
            // Lógica para redirecionar à página inicial após o recarregamento
           		window.addEventListener('beforeunload', function() {
                window.location.href = 'pagina-inicial.html'; // Substitua pela sua URL
            });
        });
		</script>
		</html>
		`)

	} else {
		templ.Execute(w, nil)
	}
}

func executarArquivoShell(w http.ResponseWriter, r *http.Request) {
	//Executar o arquivo shell
	cmd := exec.Command("./script.sh")

	//Capturar a saída da execução
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao executar o arquivo shell: %s", err), http.StatusInternalServerError)
		return
	}

	// Exibir a saída da execução
	fmt.Fprintf(w, "Saída do arquivo shell:\n%s", output)
}
