# PiteBypass

*Uma ferramenta poderosa para contornar códigos de status HTTP 4xx.*

---

## Descrição

O **PiteBypass** é uma ferramenta desenvolvida por [@pitecozz](https://github.com/pitecozz) para ajudar a contornar códigos de status HTTP 4xx (como 403, 404, etc.) em servidores web. Ele testa várias técnicas, como manipulação de verbos HTTP, adição de cabeçalhos, uso de diferentes User-Agents, extensões de arquivo, credenciais padrão, entre outras.

---

## Instalação

Siga os passos abaixo para instalar e configurar o PiteBypass.

### Pré-requisitos
- **Go** (versão 1.16 ou superior) instalado na máquina.
- **Git** instalado (para clonar o repositório).

### 1. Passos para Instalação

1.1 **Clone o Repositório**:
   ```bash
   git clone https://github.com/pitecozz/PiteBypass.git
   cd PiteBypass
   ```
1.2 **Instale as Dependências**:
   ```bash
   sudo mkdir -p /usr/local/share/pitebypass/templates
   sudo cp templates/* /usr/local/share/pitebypass/templates/
   ```
1.3 **Compile o projeto**:
   ```bash
   go build -o pitebypass
   ```
1.4 **Mova o Bin[ario para o PATH**:
   ```bash
   sudo mv pitebypass /usr/local/bin/
```
1.5 **Verifique a instalação:**
```bash
pitebypass -h
```

## Liçenca
Este projeto está licenciado sob a MIT License. Consulte o arquivo LICENSE para mais detalhes.
