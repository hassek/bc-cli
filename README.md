# Butler Coffee CLI

[![Latest Release](https://img.shields.io/github/v/release/butlercoffee/bc-cli?label=version)](https://github.com/butlercoffee/bc-cli/releases)
[![Go Version](https://img.shields.io/badge/go-1.25.4-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Built with Bubble Tea](https://img.shields.io/badge/Built%20with-Bubble%20Tea-FF69B4)](https://github.com/charmbracelet/bubbletea)

```
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMWX0kkkO0XWMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMNKko:'......'cxXMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMWN0xl;..............:0WMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMN0xl;.........''''......;0WMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMW0o;............',;;;;,.....:0WMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMNd'...............',;;;;,.....:0MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMWx..................';;;;;,.....:0MMMMMMMWNNWWMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMWo...................';;,,'......:KMWNKko:;;:lkXMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMWx.........................':l,...;dl;'..''''..cKMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMNo...................';:dOKX0c......',;;;;;;,..dWMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMNo..............':lx0KK0xl;....',;;;,,',;;;'.'kWMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMNo.........';ldkkxoc;,....',;;;,'..........;kNMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMXo.....;cdxxoc;'.....',,;,,'.....;lc,..'oONMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMXl....:l:,'....',,,,,''.....,:ldkO0ko;'lKWMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMXc.........'',''......';codxkkkkkkkkxc',kNMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMWX0xc'.......''.......';coxkkkkkkkkkkkkkkxl''dNMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMNkc,...............,:ldxxkOOOOkkOkOOkkOOxl;,,...xWMMMMMMWNNNNWMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMKc..............';ldkOOOOOOOOOOkoc::lxOOkc.'oxl'.,kXXK0kdl::::oKWMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMk...........';codxkOOOOOOOOOOOOl.....,dOk:.:OKKko:,;;;;:coxkkc.cXMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMO,..........lxxxxkO00000000000Oc.....'d0Ol.,kKkld0OkkO0KKKKKOc.cXMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMWOc'....',..lxxxxk0000000000000kl;,,:dO0Oc.;kK0kk0KKKK0KK0Od:';OWMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMWKkxxOKKc.ckxxxk0K00000000000K00OO0000k;.lOOkkkO000Okxoc,';dXWMMMMMMMMMM
MMMMWXNMMMMMMMMMMMMMMMMMMMMMMMMMMMMWk',xOxxxO000KKKKKKK0KK00KKKKK0o..cc,''',;;;,,,,cod0WNXNMMMMMMMMM
MMMXo,cKMMMMMMMMMMMMMMMMMMMMMMMMMMMMK;.l0Okxxk00KKKKKKKKKK00KKKKK0kc,',;;..coodxk0XWXlcK0:oNMMMMMMMM
MMNo...;OWMMMMMMMMMMMMMMMMMMMMMMMMMKl..;oxkOOkkO00K000KK0KKK00K0KKK0Okkd,'dNMMMMMMMMK:;Ok,lNMMMMMMMM
MWx.....'oKWMMMMMMMMMMMMMMMMMMMMMXx,......';lxO000KK00KK0KKKKKK0KKKK0kl',kNMMMMMMMMNo,xO;:KMMMMMMMMM
MK:.......'lOXWMMMMMMMMMMMMMMMN0o,............,;:clodxkOOdlccok00K0xc'..oNMMMMMMMMMNdc0KcoNMMMMMMMMM
Mk...........,cxOKXNWWWWNXK0xl;.....................;;;,'.....,oOOl......dNMMMMMWXOxocllclxOKNMMMMMM
No....,,..........,;::::;,'.......,:'...............ckKk;.......,,.......lNMMMWOc,;:cllolc:;,:kNMMMM
Nl....,,..........................cOd'...............,xKl.......''.......xWMMMXc.'lxO0000Oko,..lkKWM
Nl....,,...........,do'............lOx;................,'.......,,.......oNMMMX:...............:;,oX
Wx....',..........'dx;..............;dkxl;'...................cOXKx;....,OMMMMNl..............:K0;'k
MO,...',..........:kc.................,cdkkdl:,..........'okkKNMMMMN0ko',0MMMMMk'.........'...:xl':K
MNl....,'.........:x:.....................;coxkxl,........cXMMMW0dxXMM0,'kMMMMNk,..............':xXW
MMO,...',.........'oo'........................,ckkxc.......xWMMNo.,kWMO,,0MMXd;''............'..cOWM
MMWd....''.........'lo;.........................'oOOo......:KMMMX0KWMWx.'oxd:..:xxolcc::ccloxko,.cXM
MMMNl....''..........;d:.........................';dOc.....'OMMXdcdXMXc.........,clodxxxxxdol:;;l0WM
MMMMXl...............,dc...........................cOo......xWMKl,lKWx'...........'clcccccllox0XWMMM
MMMMMXl..............'dl..........................'lOl......dWMMWNWW0;............cXMWWWWMMMMMMMMMMM
MMMMMMNd'.............ld;........................;oxd'......xWMMMMW0:..........'..xWMMMMMMMMMMMMMMMM
MMMMMMMWO;.............cdc'.....................:xOd,......'kNNNNNO;..........'..lNMMMMMMMMMMMMMMMMM
MMMMMMMMMXo'............'cllc::;;;'...........,oxc,........;kK00Oo'............'oXMMMMMMMMMMMMMMMMMM
MMMMMMMMMMW0c..............',;;;,'.........':odc'..........ckkxo;............'cOWMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMWOl'.......................,;:llc;............,ooc;'..........,cd0WMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMWKd:....................',,'...........'....;;'.;d0KOkxxxxk0XWMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMN0d:'.....................................,ckXMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMWKkl;'..............................;lkXWMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMNKOdl:,'.................,;:ldkKNMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMWNK0OkxddooooddxkO0KXNWMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
```

> A delightful command-line interface for discovering, learning about, and ordering specialty coffee from [Butler Coffee](https://butler.coffee) - all from your terminal.

---

## Table of Contents

- [Why Butler Coffee CLI?](#why-butler-coffee-cli)
- [Installation](#installation)
  - [Via Homebrew (macOS)](#via-homebrew-macos)
  - [Download Pre-compiled Binary](#download-pre-compiled-binary)
  - [Build from Source](#build-from-source)
- [Getting Started](#getting-started)
  - [First Time Setup](#first-time-setup)
  - [Available Commands](#available-commands)
- [Learn About Coffee](#learn-about-coffee)
- [Coffee Subscriptions](#coffee-subscriptions)
- [Features](#features)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

---

## Why Butler Coffee CLI?

Butler Coffee CLI brings the entire specialty coffee experience directly to your terminal. Whether you're a developer who lives in the command line or a coffee enthusiast looking for a unique way to explore coffee, this CLI offers:

- **Interactive Learning**: Browse a comprehensive coffee knowledge base with articles on brewing methods, grinders, roasting, and more - all from your terminal
- **Coffee Subscriptions**: Configure and manage coffee subscriptions with an intuitive TUI (Terminal User Interface)
- **Complete Control**: Pause, resume, update, or cancel subscriptions without leaving your workflow
- **Delightful Experience**: Animated duck mascot and smooth Bubble Tea-powered interface make it a joy to use
- **Developer-Friendly**: Built with Go, designed for developers by developers
- **No Context Switching**: Order coffee, manage subscriptions, and learn about coffee without opening a browser

---

## What you are getting into

This is a CLI tool to find, learn, and order coffee from butler.coffee directly from your terminal! Whether you're a coffee enthusiast looking to deepen your knowledge or someone seeking the perfect coffee subscription, Butler Coffee CLI brings the entire coffee experience to your command line.

## Installation

### Via Homebrew (macOS)
```bash
# Add the Butler Coffee tap
brew tap butlercoffee/tap

# Install bc-cli
brew install bc-cli

# Or install directly in one command
brew install butlercoffee/tap/bc-cli
```

### Download Pre-compiled Binary
Download the latest release for your platform from [GitHub Releases](https://github.com/butlercoffee/bc-cli/releases)

### Build from Source
```bash
# Clone the repository
git clone https://github.com/butlercoffee/bc-cli.git
cd bc-cli

# Build the binary
make compile

# Or build directly with go
go build -o bc-cli

# Optionally, move to your PATH
sudo mv bc-cli /usr/local/bin/
```

## Getting Started

### First Time Setup

```bash
# Create a new account
bc-cli signup

# Or login with existing credentials
bc-cli login
```

### Available Commands

```bash
# Authentication
bc-cli login                # Login to your Butler Coffee account
bc-cli signup               # Create a new account
bc-cli logout               # Logout and clear stored credentials

# Learning & Discovery
bc-cli learn                # Browse coffee knowledge base interactively
bc-cli learn bookmarks      # View your saved articles (requires login)

# Shopping
bc-cli subscriptions        # Browse and subscribe to coffee subscriptions
bc-cli products             # Browse and purchase one-time coffee products

# Subscription Management (requires login)
bc-cli manage               # Manage your active subscriptions
                           # - Pause/resume subscriptions
                           # - Update quantity and preferences
                           # - View subscription details
                           # - Cancel subscriptions
```

## Learn About Coffee

Dive deep into the world of coffee with our comprehensive, interactive knowledge base:

- **Coffee Basics**: Understand the fundamentals of coffee, from bean to cup
- **Grinders**: Learn about different grinder types and how they impact your brew
- **Brewing Methods**: Explore various brewing techniques (pour-over, espresso, French press, and more)
- **Coffee Types**: Discover different coffee varieties, origins, and flavor profiles
- **Water**: Understand how water quality affects your coffee
- **Roasting**: Learn about roast levels and how they transform coffee beans

Browse articles interactively with `bc-cli learn`, save your favorite articles with bookmarks, and access them anytime with `bc-cli learn bookmarks` (requires login).

Note: We are constantly adding more content!

## Coffee Subscriptions

Butler Coffee offers three distinct subscription tiers to match your coffee journey:

### Explorer Tier
Our foundational tier delivers carefully curated, high-quality coffee selections. Perfect for those who appreciate great coffee without breaking the bank. Each shipment brings you exceptional beans that have been thoughtfully selected to expand your palate.

### Alpine Tier
Step up to our next level tier featuring rare and exclusive coffee selections from renowned origins. This subscription brings you limited-edition beans, micro-lot coffees, and unique varietals that aren't available in our standard offerings. Ideal for enthusiasts who want to explore the finer side of specialty coffee.

### "I don't care how much it costs, just give me the best of the best"
The ultimate coffee experience. This tier features the absolute finest coffees in the world—competition-winning beans, ultra-rare micro-lots, and exclusive releases that money can rarely buy. If you demand nothing but the absolute best and price is no object, this is your subscription.

This tier may or may not be available depending on the season.

## Features

- **User Authentication**: Secure login and logout with automatic token refresh
- **Account Creation**: Create a new Butler Coffee account directly from the CLI
- **Coffee Subscriptions**: Browse and subscribe to subscription tiers with interactive configuration
  - Choose your preferred grind type (whole bean or ground)
  - Select brewing method (Espresso, V60, French Press, Pour Over, Drip, Cold Brew, Moka Pot)
  - Configure monthly quantity and preferences
  - Integrated Stripe checkout with automatic browser opening
- **Subscription Management**: Comprehensive control over your active subscriptions
  - Pause and resume subscriptions at any time
  - Update quantity and coffee preferences
  - Cancel subscriptions with safety confirmations
  - View subscription details, billing info, and next shipment dates
- **Product Purchases**: Browse and purchase one-time coffee products
  - Interactive product selection with detailed information
  - Customizable grind and brewing preferences per order
  - Seamless checkout experience
- **Interactive Learning**: Access comprehensive coffee knowledge base
  - Browse by category or section
  - Read full articles with markdown formatting
  - Bookmark articles for later (authenticated users)
  - View all saved bookmarks
- **Animated TUI**: Delightful terminal interface powered by Bubble Tea
  - Animated duck mascot based on the Butler Coffee logo
  - Smooth cursor navigation and scrolling
  - Color-coded status indicators and confirmations

## Configuration

The CLI stores configuration and authentication tokens in `~/.butler-coffee/config.json`.

### Environment Variables

- **`BASE_HOSTNAME`**: Override the API URL (defaults to `https://api.butler.coffee`)
  ```bash
  export BASE_HOSTNAME=http://localhost:8000  # For local development
  ```

### Customizable Settings

You can edit `~/.butler-coffee/config.json` to customize:

- **`min_quantity`**: Minimum quantity per month for subscriptions (default: 1)
- **`max_quantity`**: Maximum quantity per month for subscriptions (default: 10)

### Supported Platforms

- macOS (uses `open` for browser integration)
- Linux (uses `xdg-open` for browser integration)
- Windows (uses `rundll32` for browser integration)

---

## Contributing

We welcome contributions! Whether it's bug reports, feature requests, or code contributions, we appreciate your help making Butler Coffee CLI better.

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute.

### Development

```bash
# Clone the repository
git clone https://github.com/butlercoffee/bc-cli.git
cd bc-cli

# Install dependencies and setup pre-commit hooks
make install

# Build the project
make compile

# Run tests
go test ./...
```

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Keywords & Topics

**For developers**: `cli`, `terminal`, `tui`, `golang`, `bubble-tea`, `command-line`, `coffee`, `cobra`

**For coffee enthusiasts**: `specialty-coffee`, `coffee-subscription`, `coffee-learning`, `coffee-ordering`, `artisan-coffee`

**Categories**: `productivity`, `terminal-app`, `interactive-cli`, `developer-tools`, `coffee-tech`

---

## Related Links

- **Website**: [butler.coffee](https://butler.coffee)
- **Issues**: [GitHub Issues](https://github.com/butlercoffee/bc-cli/issues)
- **Releases**: [GitHub Releases](https://github.com/butlercoffee/bc-cli/releases)

---

<div align="center">

Made with ☕️ and 🦆 by the Butler Coffee team

</div>
