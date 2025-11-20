Butler coffee cli tool
```
                                             ╆┨┭┫┱┭┭┭┭┣┩┨╀
                                            ┬┦┾         ╇┬┢┾
                                           ┤┩             ┾┖┹
                                          ┲┢╉     ╄┘┢┘┾   ╃┨┐
                                          ┬┬      ╅├┏┖╀ ╈┐┠┥┕┤╅
                                          ┬┢          ╃┰┓╀ ┰╄┾┤┮┰┰┱╇
                                          ╅┕┮        ╊┕█┅├┭┶╂╄╄┼┲┨┙│┸
                                           ╂┤┚┷ ╋╈╈┵┮┨┓███████████┑┷╊
                                             ╉┏┼╀┠┖┏┉▀███┟┱┱┱┱┱┱┱╅╋
                                            ╊╁╀┃┺╈   ╈╇┱▌▀┕╁
                            ╆┱┱┱┵╉     ╅┺┩┥┪┤┖┚┱┱┰      ┿┒██┨╈
                          ┼┨┨┚┾ ┸┨┨┬┨┨┨┭┪┭┪┬┴┴┴┴┴┞┪       ┬┇█┋
                          ┵┗╉╀┯┣┷    ╈┬┦┷╈       ╋┴┎┼      ┬▄█┌
                          ╉┭┛╃ ╂╉   ┹┧┽           ╆┗┴      ┫▀█┊
                            ┕┫ ┮╃╊  ├┡╀           ┳┕╄     ╀┍██┐
                            ╋┞┚┚┄┅┖┢┴┎┣┭┭┭┭┺   ╅┧├┠╈     ┺┄██┆╇
                              ┶┋███▀▀▀▄┍┏┞├├├┊┍━┠┪╋  ╉╀┦┏▀██┎┼
                                ┰┇███┊┊┇┕┗┘┚┚┬┼┿┴┴┮┥┠┈████┍┭╊
                                  ┺┠▐▌██▄┛┞┕┝┠┌▀████████┏┮╊
                                     ╇┱┱┭┓┒┒██┄┒┒│██└┞┤┼┼
                                           ┯█┃┎├┝┅▌▐███▀▌┊╉
                                           ┺┒┘┫┣┯┡┦┖┗┟┵┷
                                             ╅╀╀╀╀╀╀╊
```

# What you are getting into

This is a CLI tool to find, learn, and order coffee from butler.coffee directly from your terminal! Whether you're a coffee enthusiast looking to deepen your knowledge or someone seeking the perfect coffee subscription, Butler Coffee CLI brings the entire coffee experience to your command line.

## Learn About Coffee

Dive deep into the world of coffee with comprehensive information about:

- **Coffee Basics**: Understand the fundamentals of coffee, from bean to cup
- **Grinders**: Learn about different grinder types and how they impact your brew
- **Brewing Methods**: Explore various brewing techniques (pour-over, espresso, French press, and more)
- **Coffee Types**: Discover different coffee varieties, origins, and flavor profiles
- **Water**: Understand how water quality affects your coffee
- **Roasting**: Learn about roast levels and how they transform coffee beans

## Coffee Subscriptions

Butler Coffee offers three distinct subscription tiers to match your coffee journey:

### Butler Coffee
Our foundational tier delivers carefully curated, high-quality coffee selections. Perfect for those who appreciate great coffee without breaking the bank. Each shipment brings you exceptional beans that have been thoughtfully selected to expand your palate.

### Collection Coffee
Step up to our premium tier featuring rare and exclusive coffee selections from renowned origins. This subscription brings you limited-edition beans, micro-lot coffees, and unique varietals that aren't available in our standard offerings. Ideal for enthusiasts who want to explore the finer side of specialty coffee.

### "I don't care how much it costs, just give me the best of the best"
The ultimate coffee experience. This tier features the absolute finest coffees in the world—competition-winning beans, ultra-rare micro-lots, and exclusive releases that money can rarely buy. If you demand nothing but the absolute best and price is no object, this is your subscription.

### Daily Coffee Snippets
Subscribe to receive a daily coffee snippet delivered straight to your email. Each day, you'll get bite-sized knowledge about coffee—from brewing tips to origin stories, tasting notes to industry insights. Perfect for coffee lovers who want to learn something new every day.

### Jura Coffee Machines
We exclusively offer Jura coffee machines—the Swiss-engineered gold standard in automatic coffee makers. Browse our selection of Jura machines perfect for home use, from entry-level models to professional-grade equipment. Each machine delivers barista-quality espresso, coffee, and milk-based drinks at the touch of a button.

**Office Solutions**: Need to elevate your workplace coffee experience? We arrange complete Jura coffee solutions for offices of any size, including machine selection, installation, maintenance, and ongoing coffee supply.

## Features

- **User Authentication**: Secure login to your Butler Coffee account using OAuth2
- **Account Creation**: Create a new Butler Coffee account directly from the CLI
- **Purchase & Subscribe**: Browse and purchase coffee subscriptions through our API integration
- **Interactive Learning**: Access our comprehensive coffee knowledge base
- **Order Management**: Track and manage your coffee subscriptions
- **Curated Tools & Gear**: Discover coffee tools and equipment we personally use and recommend

## Installation

### Via Homebrew
```bash
brew tap butlercoffee/tap
brew install butler-coffee
```

### Download Pre-compiled Binary
Download the latest release for your platform from [GitHub Releases](https://github.com/butlercoffee/bc-cli/releases)

### Build from Source
```bash
# Clone the repository
git clone https://github.com/butlercoffee/bc-cli.git
cd bc-cli

# Build the binary
go build -o bc

# Optionally, move to your PATH
sudo mv bc /usr/local/bin/
```

## Getting Started

```bash
# Login to your account
bc login

# Or create a new account
bc signup

# Explore coffee knowledge
bc learn

# Browse subscriptions
bc subscriptions

# Subscribe to daily snippets
bc subscribe-snippets

# Discover recommended coffee tools
bc tools
```
