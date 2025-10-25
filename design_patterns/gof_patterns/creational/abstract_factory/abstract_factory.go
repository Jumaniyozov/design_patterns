// Package abstract_factory demonstrates the Abstract Factory pattern in Go.
// Abstract Factory provides an interface for creating families of related
// or dependent objects without specifying their concrete classes.
package abstract_factory

import "fmt"

// Button interface represents an abstract button product.
type Button interface {
	Render()
	OnClick()
}

// Checkbox interface represents an abstract checkbox product.
type Checkbox interface {
	Render()
	Toggle()
}

// TextField interface represents an abstract text field product.
type TextField interface {
	Render()
	SetText(text string)
}

// GUIFactory is the abstract factory interface.
type GUIFactory interface {
	CreateButton() Button
	CreateCheckbox() Checkbox
	CreateTextField() TextField
	GetName() string
}

// Windows implementations
type WindowsButton struct{}

func (b *WindowsButton) Render() {
	fmt.Println("[Windows] Rendering rectangular button with gradient")
}

func (b *WindowsButton) OnClick() {
	fmt.Println("[Windows] Button clicked with Windows animation")
}

type WindowsCheckbox struct{}

func (c *WindowsCheckbox) Render() {
	fmt.Println("[Windows] Rendering square checkbox")
}

func (c *WindowsCheckbox) Toggle() {
	fmt.Println("[Windows] Checkbox toggled with checkmark animation")
}

type WindowsTextField struct {
	text string
}

func (t *WindowsTextField) Render() {
	fmt.Printf("[Windows] Rendering text field with blue border: '%s'\n", t.text)
}

func (t *WindowsTextField) SetText(text string) {
	t.text = text
}

// WindowsFactory creates Windows UI components.
type WindowsFactory struct{}

func (f *WindowsFactory) CreateButton() Button {
	return &WindowsButton{}
}

func (f *WindowsFactory) CreateCheckbox() Checkbox {
	return &WindowsCheckbox{}
}

func (f *WindowsFactory) CreateTextField() TextField {
	return &WindowsTextField{}
}

func (f *WindowsFactory) GetName() string {
	return "Windows"
}

// Mac implementations
type MacButton struct{}

func (b *MacButton) Render() {
	fmt.Println("[Mac] Rendering rounded button with sleek design")
}

func (b *MacButton) OnClick() {
	fmt.Println("[Mac] Button clicked with smooth Mac transition")
}

type MacCheckbox struct{}

func (c *MacCheckbox) Render() {
	fmt.Println("[Mac] Rendering rounded checkbox")
}

func (c *MacCheckbox) Toggle() {
	fmt.Println("[Mac] Checkbox toggled with fade animation")
}

type MacTextField struct {
	text string
}

func (t *MacTextField) Render() {
	fmt.Printf("[Mac] Rendering text field with rounded corners: '%s'\n", t.text)
}

func (t *MacTextField) SetText(text string) {
	t.text = text
}

// MacFactory creates Mac UI components.
type MacFactory struct{}

func (f *MacFactory) CreateButton() Button {
	return &MacButton{}
}

func (f *MacFactory) CreateCheckbox() Checkbox {
	return &MacCheckbox{}
}

func (f *MacFactory) CreateTextField() TextField {
	return &MacTextField{}
}

func (f *MacFactory) GetName() string {
	return "Mac"
}

// Linux implementations
type LinuxButton struct{}

func (b *LinuxButton) Render() {
	fmt.Println("[Linux] Rendering flat button with GTK theme")
}

func (b *LinuxButton) OnClick() {
	fmt.Println("[Linux] Button clicked with ripple effect")
}

type LinuxCheckbox struct{}

func (c *LinuxCheckbox) Render() {
	fmt.Println("[Linux] Rendering checkbox with theme colors")
}

func (c *LinuxCheckbox) Toggle() {
	fmt.Println("[Linux] Checkbox toggled")
}

type LinuxTextField struct {
	text string
}

func (t *LinuxTextField) Render() {
	fmt.Printf("[Linux] Rendering GTK text field: '%s'\n", t.text)
}

func (t *LinuxTextField) SetText(text string) {
	t.text = text
}

// LinuxFactory creates Linux UI components.
type LinuxFactory struct{}

func (f *LinuxFactory) CreateButton() Button {
	return &LinuxButton{}
}

func (f *LinuxFactory) CreateCheckbox() Checkbox {
	return &LinuxCheckbox{}
}

func (f *LinuxFactory) CreateTextField() TextField {
	return &LinuxTextField{}
}

func (f *LinuxFactory) GetName() string {
	return "Linux"
}

// Database example
type Connection interface {
	Connect() error
	Close() error
}

type QueryBuilder interface {
	Select(fields ...string) QueryBuilder
	From(table string) QueryBuilder
	Build() string
}

type Transaction interface {
	Begin() error
	Commit() error
	Rollback() error
}

type DatabaseFactory interface {
	CreateConnection() Connection
	CreateQueryBuilder() QueryBuilder
	CreateTransaction() Transaction
	GetDBType() string
}

// MySQL implementations
type MySQLConnection struct{}

func (c *MySQLConnection) Connect() error {
	fmt.Println("[MySQL] Connected to MySQL database")
	return nil
}

func (c *MySQLConnection) Close() error {
	fmt.Println("[MySQL] Closed MySQL connection")
	return nil
}

type MySQLQueryBuilder struct {
	query string
}

func (q *MySQLQueryBuilder) Select(fields ...string) QueryBuilder {
	q.query = "SELECT " + fields[0]
	return q
}

func (q *MySQLQueryBuilder) From(table string) QueryBuilder {
	q.query += " FROM " + table
	return q
}

func (q *MySQLQueryBuilder) Build() string {
	return q.query + " /* MySQL syntax */"
}

type MySQLTransaction struct{}

func (t *MySQLTransaction) Begin() error {
	fmt.Println("[MySQL] BEGIN TRANSACTION")
	return nil
}

func (t *MySQLTransaction) Commit() error {
	fmt.Println("[MySQL] COMMIT")
	return nil
}

func (t *MySQLTransaction) Rollback() error {
	fmt.Println("[MySQL] ROLLBACK")
	return nil
}

type MySQLFactory struct{}

func (f *MySQLFactory) CreateConnection() Connection {
	return &MySQLConnection{}
}

func (f *MySQLFactory) CreateQueryBuilder() QueryBuilder {
	return &MySQLQueryBuilder{}
}

func (f *MySQLFactory) CreateTransaction() Transaction {
	return &MySQLTransaction{}
}

func (f *MySQLFactory) GetDBType() string {
	return "MySQL"
}

// PostgreSQL implementations
type PostgreSQLConnection struct{}

func (c *PostgreSQLConnection) Connect() error {
	fmt.Println("[PostgreSQL] Connected to PostgreSQL database")
	return nil
}

func (c *PostgreSQLConnection) Close() error {
	fmt.Println("[PostgreSQL] Closed PostgreSQL connection")
	return nil
}

type PostgreSQLQueryBuilder struct {
	query string
}

func (q *PostgreSQLQueryBuilder) Select(fields ...string) QueryBuilder {
	q.query = "SELECT " + fields[0]
	return q
}

func (q *PostgreSQLQueryBuilder) From(table string) QueryBuilder {
	q.query += " FROM " + table
	return q
}

func (q *PostgreSQLQueryBuilder) Build() string {
	return q.query + " /* PostgreSQL syntax with RETURNING */"
}

type PostgreSQLTransaction struct{}

func (t *PostgreSQLTransaction) Begin() error {
	fmt.Println("[PostgreSQL] BEGIN")
	return nil
}

func (t *PostgreSQLTransaction) Commit() error {
	fmt.Println("[PostgreSQL] COMMIT")
	return nil
}

func (t *PostgreSQLTransaction) Rollback() error {
	fmt.Println("[PostgreSQL] ROLLBACK")
	return nil
}

type PostgreSQLFactory struct{}

func (f *PostgreSQLFactory) CreateConnection() Connection {
	return &PostgreSQLConnection{}
}

func (f *PostgreSQLFactory) CreateQueryBuilder() QueryBuilder {
	return &PostgreSQLQueryBuilder{}
}

func (f *PostgreSQLFactory) CreateTransaction() Transaction {
	return &PostgreSQLTransaction{}
}

func (f *PostgreSQLFactory) GetDBType() string {
	return "PostgreSQL"
}

// Application uses the factories
type Application struct {
	factory GUIFactory
}

func NewApplication(factory GUIFactory) *Application {
	return &Application{factory: factory}
}

func (app *Application) RenderUI() {
	fmt.Printf("\nRendering UI with %s theme:\n", app.factory.GetName())

	button := app.factory.CreateButton()
	checkbox := app.factory.CreateCheckbox()
	textField := app.factory.CreateTextField()

	button.Render()
	checkbox.Render()
	textField.SetText("Hello, " + app.factory.GetName() + "!")
	textField.Render()
}

func (app *Application) HandleUserInput() {
	fmt.Printf("\nHandling user input with %s components:\n", app.factory.GetName())

	button := app.factory.CreateButton()
	checkbox := app.factory.CreateCheckbox()

	button.OnClick()
	checkbox.Toggle()
}
