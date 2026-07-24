package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProvider_LibraryCreateMovies(t *testing.T) {
	resource.Test(t, resource.TestCase {
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep {
			{
				Config: `
					provider "jellyfin" {
						endpoint = "http://localhost:8097"
						username = "admin"
						password = "admin"
					}

					resource "jellyfin_library" "movies" {
						name = "Movies"
						collection_type = "movies"
						paths = ["/media"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("jellyfin_library.movies", "id"),
					resource.TestCheckResourceAttr("jellyfin_library.movies", "name", "Movies"),
					resource.TestCheckResourceAttr("jellyfin_library.movies", "collection_type", "movies"),
				),
			},
		},
	})
}

func TestAccProvider_LibraryCreateTvShows(t *testing.T) {
	resource.Test(t, resource.TestCase {
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep {
			{
				Config: `
					provider "jellyfin" {
						endpoint = "http://localhost:8097"
						username = "admin"
						password = "admin"
					}

					resource "jellyfin_library" "tvshows" {
						name = "TV Shows"
						collection_type = "tvshows"
						paths = ["/media"]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("jellyfin_library.movies", "id"),
					resource.TestCheckResourceAttr("jellyfin_library.movies", "name", "TV Shows"),
					resource.TestCheckResourceAttr("jellyfin_library.movies", "collection_type", "tvshows"),
				),
			},
		},
	})
}
