// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package main

import (
	"github.com/rs/zerolog/log"
	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_sdk/cmd"
	"github.com/trustero/api/go/receptor_v1"
)

const (
	receptorName = "trr-custom"
	serviceName1 = "Translate"
)

// This struct holds the credentials the receptor needs to authenticate with the
// service provider. A display name and placeholder tag should be provided
// which will be used in the UI when activating the receptor.
// This is what will be returned in the GetCredentialObj call
type Receptor struct {
	FirstName   string `trustero:"display:Primary User;placeholder:First Name"`
}

// Defines the structure of a single row of our "Translation" evidence
type TrusteroTranslationRow struct {
    LangId           string     `trustero:"id:"`
    Language         string     `trustero:"display:Language;order:1"`
    Phrase           string     `trustero:"display:Phrase;order:2"`
}

// Set the name of the receptor in the const declaration above
// This will let the receptor inform Trustero about itself
func (r *Receptor) GetReceptorType() string {
	return receptorName
}

// Set the names of the services in the const declaration above
// This will let the receptor inform Trustero about itself
// Feel free to add or remove services as needed
func (r *Receptor) GetKnownServices() []string {
	return []string{serviceName1}
}

// This will return Receptor struct defined above when the receptor is asked to
// identify itself
func (r *Receptor) GetCredentialObj() (credentialObj interface{}) {
	return r
}

// This function will call into the service provider API with the provided
// credentials and confirm that the credentials are valid. Usually a simple
// API call like GET org name. If the credentials are not valid,
// return a relevant error message
func (r *Receptor) Verify(credentials interface{}) (ok bool, err error) {
	c := credentials.(*Receptor)
	log.Info().Msgf("verify: checking credentials for %s", c.FirstName)
	ok = (err == nil)
	return
}

// The Discover function returns a list of Service Entities. This function
// makes any relevant API calls to the Service Provider to gather information
// about how many Service Entity Instances are in use. If at any point this
// function runs into an error, log that error and continue
func (r *Receptor) Discover(credentials interface{}) (svcs []*receptor_v1.ServiceEntity, err error) {
    services := receptor_sdk.NewServiceEntities()
	services.AddService(serviceName1, "Language", "English", "en")
	services.AddService(serviceName1, "Language", "German", "de")
	services.AddService(serviceName1, "Language", "Italian", "it")

	return services.Entities, err
}

// Report will often make the same API calls made in the Discover call, but it
// will additionally create evidences with the data returned from the API calls
func (r *Receptor) Report(credentials interface{}) (evidences []*receptor_sdk.Evidence, err error) {
	c := credentials.(*Receptor)
	report := receptor_sdk.NewReport()

	caption := serviceName1 + " Hello Translations"
	description := "List of translations of the term 'Hello, [user]'"
	evidence := receptor_sdk.NewEvidence(
		serviceName1,
		"Language",
		caption,
        description)

	// Note the API call made and results returned for English
	evidence.AddSource("https://translate.google.com/?sl=auto&tl=en&text=Hello&op=translate", "Hello")
	evidence.AddRow(TrusteroTranslationRow{
		LangId:           "en",
		Language:         "English",
		Phrase:           "Hello, " + c.FirstName,
	})
	// Note the API call made and results returned for Italian
	evidence.AddSource("https://translate.google.com/?sl=auto&tl=it&text=Hello&op=translate", "Ciao")
	evidence.AddRow(TrusteroTranslationRow{
		LangId:           "it",
		Language:         "Italian",
		Phrase:           "Ciao, " + c.FirstName,
	})
	// Note the API call made and results returned for German
	evidence.AddSource("https://translate.google.com/?sl=auto&tl=de&text=Hello&op=translate", "Hallo")
	evidence.AddRow(TrusteroTranslationRow{
		LangId:           "de",
		Language:         "German",
		Phrase:           "Hallo, " + c.FirstName,
	})

	report.AddEvidence(evidence)
	return report.Evidences, err
}

func main() {
	cmd.Execute(&Receptor{})
}
