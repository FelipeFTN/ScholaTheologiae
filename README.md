# 📚 Schola Theologiae ✝️
### *Free & Open-Source Catholic Knowledge for Everyone*

> *"To one who has faith, no explanation is necessary. To one without faith, no explanation is possible."* \
— **St. Thomas Aquinas**

---

## 🕊️ About

**Schola Theologiae** is a comprehensive digital library dedicated to providing free access to the treasures of Catholic theological knowledge. Built with modern web technologies, it serves the timeless wisdom of the Church Fathers, Saints, and the Magisterium to anyone seeking deeper understanding of the Catholic faith.

### ✨ Features

- 📖 **Summa Theologiae** - Complete works of St. Thomas Aquinas
- 📜 **Catecism** - St. Pius X's comprehensive catechism
- 🎯 **Structured Navigation** - Organized by parts, questions, and articles
- 🔍 **Search Functionality** - Find specific teachings and concepts
- 📱 **Mobile-Friendly** - Responsive design for all devices
- 🌐 **Multi-language Support** - Latin and vernacular texts
- ⚡ **Fast API** - Go-powered backend with SQLite database

---

## 🏗️ Architecture

### Backend (Go API)
```
api/
├── controllers/     # Request handlers
├── data/           # Database operations
├── models/         # Data structures
├── services/       # Business logic
└── server/         # HTTP server setup
```

### Frontend (Ruby on Rails)
```
app/
├── controllers/    # Web controllers
├── views/         # HTML templates
├── helpers/       # View helpers
└── assets/        # CSS and JavaScript
```

### Database Structure
- **SQLite** databases for efficient text storage and retrieval
- Indexed by parts, questions, articles, and chapters
- Optimized for fast searching and navigation
- Everything is handled by the [scripts](scripts/) in the `scripts/` directory — these scripts are responsible for creating the database, populating it with data, and managing migrations.

---

## 🚀 Quick Start

### Prerequisites
- **Docker** and **Docker Compose**
- **Go 1.21+** (for development)
- **Ruby 3.2+** (for development)
- **Make** (build automation)

### Using Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/FelipeFTN/ScholaTheologiae.git
cd ScholaTheologiae

# Build and run with Docker
make build
```

### Local Development

```bash
# Start the API server
./run-api.sh

# Start the Rails server
./run-app.sh
```

Visit `https://scholatheologiae.com` to explore the digital library!

---

## 📋 API Endpoints

### Books
Books are addressed by their database name: `summa_theologiae`, `catecismo_pio_x`, `confissoes`, `didaque`.

- `GET /v1/books/{name}` - List the parts of a book
- `GET /v1/books/{name}/{part}` - List the chapters of a part
- `GET /v1/books/{name}/{part}/{chapter}` - Get a chapter body (markdown)

### Search
- `GET /v1/search?q={query}` - Full-text search over every chapter title and body

Accent-insensitive, ordered by relevance inside each book, and returns a snippet with the matched terms wrapped in `<mark>`. It is backed by an FTS5 index that the API build (`make build`) creates from the chapter files; books without an index fall back to a title-only match. Article-level reads (`/v1/books/{name}/{part}/{chapter}/{article}`) answer 404 while no book ships articles.

### Health Check
- `GET /v1/health` - API health status

---

## 🎨 Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Backend API** | Go + Gin | Fast, concurrent API server |
| **Frontend** | Ruby on Rails | Web application framework |
| **Database** | SQLite | Lightweight, file-based database |
| **Containerization** | Docker | Deployment and portability |
| **Process Management** | Supervisor | Multi-process container management |
| **Web Server** | Nginx | Reverse proxy and static files |
| **Build System** | Make | Build automation |

---

> I'm literally extracting the max processing capacity possible from one dyno in Heroku 😁

## 📖 Available Works

### 📚 Summa Theologiae by St. Thomas Aquinas
> *"Three things are necessary for the salvation of man: to know what he ought to believe; to know what he ought to desire; and to know what he ought to do."* - **St. Thomas Aquinas**

The greatest systematic work of theology and philosophy, containing:
- **Prima Pars** - God, Creation, and Angels
- **Prima Secundae** - Human Acts and Law
- **Secunda Secundae** - Virtues and Vices
- **Tertia Pars** - Christ and the Sacraments

### 📜 Catecismo Maior by St. Pius X
> *"It is not enough to have faith; we must also know what we believe."* - **St. Pius X**

Comprehensive catechetical instruction covering:
- **Fundamental Doctrines** - Creed and basic beliefs
- **Moral Theology** - Commandments and Christian living
- **Sacramental Life** - Seven sacraments explained
- **Prayer and Devotion** - Spiritual practices

---

## 🤝 Contributing

We welcome contributions from fellow Catholics and developers! Here's how you can help:

### Ways to Contribute
- 📝 Add new theological texts
- 🐛 Report bugs and issues
- 💡 Suggest new features
- 🌍 Improve translations
- 📚 Enhance documentation
- 🎨 Improve UI/UX design

### Development Setup
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

---

## 📜 Patristic Wisdom

> *"All that is true, by whomsoever it has been said, is from the Holy Spirit."* - **St. Thomas Aquinas**

> *"The Word of God is not a sounding but a piercing word, not pronounceable by the tongue but efficacious in the mind."* - **St. Bernard of Clairvaux**

> *"You must ask God to give you power to fight against the sin of pride which is your greatest enemy - the root of all that is evil, and the failure of all that is good."* - **St. Vincent de Paul**

---

## 🛡️ Deployment

### Heroku Deployment
The application is configured for Heroku deployment with:
- `heroku.yml` - Container deployment configuration
- `Procfile` - Process management
- Environment variables for production settings

### Docker Production
```bash
# Build production image
docker build -t schola-theologiae:production .

# Run with production configuration
docker run -p 8000:8000 \
  -e RAILS_ENV=production \
  -e SECRET_KEY_BASE=your_secret_key \
  schola-theologiae:production
```

---

## 📊 Project Status

- ✅ **Summa Theologiae** - Complete and searchable
- ✅ **Catecismo Pio X** - Complete and searchable
- 🚧 **Patristica Collection** - In development
- 🚧 **Divine Office** - Planned
- 🚧 **Daily Liturgy** - Planned
- 🚧 **Search Enhancement** - In progress

---

## 📄 License

This project is dedicated to the greater glory of God and the salvation of souls. Released under the MIT License - see [LICENSE](LICENSE) file for details.

---

## 🌟 Acknowledgments

Special thanks to:
- **St. Thomas Aquinas** - The Angelic Doctor
- **St. Pius X** - Pope of the Eucharist and Catechesis
- **The Church Fathers** - Pillars of theological wisdom
- **Open Source Community** - Making knowledge freely available

<div align="center">

> *"The soul that is in love with God is a soul that is in search of truth"* \
*Ad Majorem Dei Gloriam* ✝️

[![Docker](https://img.shields.io/badge/Docker-Ready-blue?logo=docker)](https://docker.com)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org)
[![Ruby](https://img.shields.io/badge/Ruby-3.2+-CC342D?logo=ruby)](https://ruby-lang.org)
[![Rails](https://img.shields.io/badge/Rails-7.0+-CC0000?logo=rubyonrails)](https://rubyonrails.org)

</div>
