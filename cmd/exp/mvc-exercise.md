### Ex1 - What does MVC stand for?
`models` | `view` | `controller`

### Ex2 - What is each layer of MVC responsible for?
- `models` are typically used for data, logic, rules; usually database
- `views` are typically used for rendering things in this case html
- `controllers` are typically used to connect everything. Accepts user input, passes that over to models to do stuff and then passes it to views in order to render it out.
  
MVC doesn't need to follow the naming of models, views, controllers. There's a web developement eco-system for Go called Buffalo and they use an MVC structure like `actions` which would be similar to `models`.

### Ex3 - What are some benefits and disadvantages to using MVC?
#### Benefits:
- Organises large scale applications
- Easier to modify and understand
- Faster deployment

#### Disadvantages:
- It can increase the complexity of your code
- It can be harder to understand the MVC Architecture
- Must have strict rules implemented


### Ex4 - Read about other ways to structure code
There are many other ways to structure code, for example; Flat structure, dependency based, domain driven design, and so on.