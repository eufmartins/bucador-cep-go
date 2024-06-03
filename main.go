package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Estrutura que representa o retorno da API ViaCEP
type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
}

func main() {
	// Define o handler para a rota raiz e inicia o servidor HTTP na porta 8080
	http.HandleFunc("/", BuscaCEPHandler)
	http.ListenAndServe(":8080", nil)
}

// Handler para a rota raiz
func BuscaCEPHandler(w http.ResponseWriter, r *http.Request) {
	// Verifica se a rota é a raiz, se não for, retorna 404 Not Found
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Obtém o parâmetro 'cep' da URL
	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		// Se o parâmetro 'cep' estiver vazio, retorna 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Chama a função BuscaCEP para obter os dados do CEP informado
	cep, err := BuscaCEP(cepParam)
	if err != nil {
		// Se ocorrer algum erro na busca do CEP, retorna 500 Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Define o cabeçalho da resposta como JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Codifica a estrutura 'cep' em JSON e escreve na resposta
	json.NewEncoder(w).Encode(cep)
}

// Função que busca informações de um CEP usando a API ViaCEP
func BuscaCEP(cep string) (*ViaCEP, error) {
	// Faz uma requisição GET para a API ViaCEP
	res, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		// Se ocorrer algum erro na requisição, retorna o erro
		return nil, err
	}
	defer res.Body.Close()
	// Lê o corpo da resposta
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// Se ocorrer algum erro na leitura do corpo, retorna o erro
		return nil, err
	}
	// Cria uma variável do tipo ViaCEP para armazenar os dados decodificados
	var c ViaCEP
	// Decodifica o JSON do corpo da resposta para a variável 'c'
	err = json.Unmarshal(body, &c)
	if err != nil {
		// Se ocorrer algum erro na decodificação, retorna o erro
		return nil, err
	}
	// Retorna a estrutura 'c' preenchida com os dados do CEP
	return &c, nil
}
