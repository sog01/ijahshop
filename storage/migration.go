package storage

// Migrate to migrate table into database
func (s Storage) Migrate() error {
	// create table product
	_, err := s.DB.Exec(
		`CREATE TABLE IF NOT EXISTS product (
			product_id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(30) NOT NULL,
			sku VARCHAR(30) NOT NULL UNIQUE,
			stock INT UNSIGNED NOT NULL
	)`)
	if err != nil {
		return err
	}

	// create table purchase
	_, err = s.DB.Exec(
		`CREATE TABLE IF NOT EXISTS purchase (
			purchase_id INTEGER PRIMARY KEY AUTOINCREMENT,
			product_id INT UNSIGNED NOT NULL,
			quantity_order INT UNSIGNED NOT NULL,
			quantity_accepted INT UNSIGNED NOT NULL,
			description TEXT DEFAULT (''),
			invoice_number VARCHAR(30) NOT NULL,
			cost DECIMAL(10, 2) NOT NULL,
			date TIMESTAMPS NOT NULL,	
			is_finish BOOLEAN DEFAULT (0)					
	)`)
	if err != nil {
		return err
	}

	// create table purchase detail
	_, err = s.DB.Exec(
		`CREATE TABLE IF NOT EXISTS purchase_detail (
			purchase_detail_id INTEGER PRIMARY KEY AUTOINCREMENT,
			purchase_id INT UNSIGNED NOT NULL,
			quantity INT UNSIGNED NOT NULL,
			description TEXT DEFAULT (''),			
			date TIMESTAMPS NOT NULL						
	)`)
	if err != nil {
		return err
	}

	// create table orders
	_, err = s.DB.Exec(
		`CREATE TABLE IF NOT EXISTS orders (
			order_id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id_format VARCHAR(30) NOT NULL,
			product_id INT UNSIGNED NOT NULL,
			quantity INT UNSIGNED NOT NULL,
			description TEXT NOT NULL,			
			date TIMESTAMPS NOT NULL,
			price DECIMAL(10, 2) NOT NULL					
	)`)
	if err != nil {
		return err
	}

	return nil
}

// Rollback to rollback all table from database
func (s Storage) Rollback() error {
	// drop table product
	_, err := s.DB.Exec("DROP TABLE product")
	if err != nil {
		return err
	}

	// drop table purchase
	_, err = s.DB.Exec("DROP TABLE purchase")
	if err != nil {
		return err
	}

	// drop table purchase detail
	_, err = s.DB.Exec("DROP TABLE purchase_detail")
	if err != nil {
		return err
	}

	// drop table orders
	_, err = s.DB.Exec("DROP TABLE orders")
	if err != nil {
		return err
	}

	return nil
}
