# 🌤️ CEP-Clima: Consulta de Clima por CEP

Um sistema em Go que consulta o clima atual de uma cidade com base no CEP (Código de Endereçamento Postal) brasileiro. A aplicação integra as APIs [ViaCEP](https://viacep.com.br/) e [WeatherAPI](https://www.weatherapi.com/) para fornecer a temperatura em **Celsius, Fahrenheit e Kelvin**.

A aplicação pode ser executada localmente com Docker ou implantada no **Google Cloud Run**.

---

## ✨ Funcionalidades

- 🔍 Consulta de clima por CEP (8 dígitos)
- 📍 Identifica automaticamente a cidade com o ViaCEP
- 🌡️ Exibe temperatura em Celsius, Fahrenheit e Kelvin
- 🖥️ Interface web responsiva para busca
- ✅ Tratamento de erros conforme especificado:
  - `422` para CEP inválido
  - `404` para CEP não encontrado
- 🐳 Execução com Docker e Docker Compose
- ☁️ Pronto para deploy no Google Cloud Run

---

## 🚀 Como Usar

### 1. Execução com Docker (Recomendado)

Certifique-se de ter o **Docker** e o **Docker Compose** instalados.

```bash
# Clone o repositório
git clone https://github.com/seu-usuario/cep-clima.git
cd cep-clima

# Crie um arquivo .env com sua chave da WeatherAPI
echo "WEATHER_API_KEY=sua_chave_aqui" > .env

# Inicie a aplicação
docker-compose up --build

# Acesse a página principal do projeto a URL_BASE e informe o CEP que deseja consultar.
# local
localhost:8080

# Cloud Run
https://cep-clima-2fcsjepn5q-uc.a.run.app/
