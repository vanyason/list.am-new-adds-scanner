package postgres

import "testing"

const (
	tableNameAll        = "ads_all"
	tableNameNew        = "ads_new"
	tableNameDoNotExist = "do_not_exist"
	tableNameToDrop     = "to_drop"
	userName            = "test"
	password            = "test"
	port                = "5432"
	dbName              = "test"
)

func defaultDriver(t *testing.T) Postgres {
	p, err := New(userName, password, port, dbName)
	if err != nil {
		t.Fatalf("Can not create db connection to test : %s", err)
	}
	return p
}

// Test default creation
func TestNewOk(t *testing.T) {
	p := defaultDriver(t)
	p.Close()
}

// Test wrong credentials
func TestNewWrongCredentials(t *testing.T) {
	p, err := New("bla", "bla", "bla", "bla")
	if err == nil {
		t.Fatalf("Created connection with wrong Credentials : %s", err)
	}

	p.Close()
}

// Test helper func to prevent sql injection
func TestCheckStringsForSql(t *testing.T) {
	const (
		ok1    = "aBcd"
		ok2    = "abBcd___dcSba"
		notOk1 = "1ab"
		notOk2 = "?2a"
	)

	if !isLegal(ok1) {
		t.Errorf("check string failed for %s", ok1)
	}

	if !isLegal(ok2) {
		t.Errorf("check string failed for %s", ok2)
	}

	if isLegal(notOk1) {
		t.Errorf("check string failed for %s", notOk1)
	}

	if isLegal(notOk2) {
		t.Errorf("check string failed for %s", notOk2)
	}
}

// Test ads table creation
func TestCreateAdsTableOk(t *testing.T) {
	p := defaultDriver(t)
	defer p.Close()

	err := p.CreateAdsTable(tableNameAll)
	if err != nil {
		t.Fatalf("Can not create table : %s", err)
	}
}

// Test table existence
func TestExistAdsTable(t *testing.T) {
	p := defaultDriver(t)
	defer p.Close()

	err := p.CreateAdsTable(tableNameAll)
	if err != nil {
		t.Fatalf("Can not create table : %s", err)
	}

	exists, err := p.ExistAdsTable(tableNameAll)
	if err != nil {
		t.Fatalf("Error while checking table existence : %s", err)
	}

	if !exists {
		t.Fatalf("Table %s doesn`t exist but should", tableNameAll)
	}
}

// Test table existence
func TestNotExistAdsTable(t *testing.T) {
	p := defaultDriver(t)
	defer p.Close()

	exists, err := p.ExistAdsTable(tableNameDoNotExist)
	if err != nil {
		t.Fatalf("Error while checking table existence : %s", err)
	}

	if exists {
		t.Fatalf("Table %s should not exist", tableNameDoNotExist)
	}
}

// Test drop table
func TestDropTable(t *testing.T) {
	p := defaultDriver(t)
	defer p.Close()

	if err := p.CreateAdsTable(tableNameToDrop); err != nil {
		t.Fatalf("Can not create table : %s", err)
	}

	exists, err := p.ExistAdsTable(tableNameToDrop)
	if err != nil {
		t.Fatalf("Error while checking table existence : %s", err)
	}
	if !exists {
		t.Fatalf("Table \"%s\" does not exist but should", tableNameToDrop)
	}

	if err := p.DropTable(tableNameToDrop); err != nil {
		t.Fatalf("Can not drop table \"%s\"", err)
	}

	exists, err = p.ExistAdsTable(tableNameToDrop)
	if err != nil {
		t.Fatalf("Error while checking table existence : %s", err)
	}
	if exists {
		t.Fatalf("Table \"%s\" exists but should", tableNameToDrop)
	}
}
