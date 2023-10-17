package cassandra

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	testAccProviders map[string]*schema.Provider
	testAccProvider  *schema.Provider
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"cassandra": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func TestProvider_configure1(t *testing.T) {
	rc := terraform.NewResourceConfigRaw(map[string]interface{}{
		"username": "cassanrda",
		"password": "cassanrda",
		"port":     9042,
		"host":     "asdf",
	})
	p := Provider()
	err := p.Configure(context.Background(), rc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProvider_configure2(t *testing.T) {
	rc := terraform.NewResourceConfigRaw(map[string]interface{}{
		"username": "cassanrda",
		"password": "cassanrda",
		"port":     9042,
		"hosts":    []interface{}{"asd"},
	})
	p := Provider()
	err := p.Configure(context.Background(), rc)
	if err != nil {
		t.Fatal(err)
	}
}

func testAccPreCheck(t *testing.T) {
	url := os.Getenv("CASSANDRA_HOST")
	if url == "" {
		t.Fatal("CASSANDRA_HOST must be set for acceptance tests")
	}

	err := testAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
	if err != nil {
		t.Fatal(err)
	}
}
