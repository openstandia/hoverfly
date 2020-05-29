package hoverfly_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	functional_tests "github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
	"github.com/dghubble/sling"
)

var _ = Describe("When running Hoverfly as a webserver", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("and its in simulate mode with flag -journal-file=journal.log", func() {

		BeforeEach(func() {
			hoverfly.Start("-webserver", "-journal-file=journal.log")
			hoverfly.SetMode("simulate")
			hoverfly.ImportSimulation(testdata.JsonGetAndPost)
		})

		Context("I can request an endpoint", func() {
			Context("using GET", func() {
				It("and it should write journal message (HTTP request/response) into file journal.log", func() {
					request := sling.New().Get("http://localhost:" + hoverfly.GetProxyPort() + "/path1")
					response := functional_tests.DoRequest(request)

					responseBody, err := ioutil.ReadAll(response.Body)
					Expect(err).To(BeNil())

					Expect(string(responseBody)).To(Equal("body1"))

					journal, err := hoverfly.GetLogFile("journal.log")
					Expect(err).To(BeNil())
					Expect(journal).To(ContainSubstring(`"path":"/path1"`))
					Expect(journal).To(ContainSubstring(`"scheme":"http"`))
					Expect(journal).To(ContainSubstring(`"method":"GET"`))
					Expect(journal).To(ContainSubstring(`"body":"` + string(responseBody) + `"`))
					Expect(journal).To(ContainSubstring(`"status":201`))
				})
			})

			Context("using POST", func() {
				It("and it should write journal message (HTTP request/response) into file journal.log", func() {
					request := sling.New().Post("http://localhost:" + hoverfly.GetProxyPort() + "/path2/resource")

					response := functional_tests.DoRequest(request)

					responseBody, err := ioutil.ReadAll(response.Body)
					Expect(err).To(BeNil())

					Expect(string(responseBody)).To(Equal("POST body response"))

					journal, err := hoverfly.GetLogFile("journal.log")
					Expect(err).To(BeNil())
					Expect(journal).To(ContainSubstring(`"path":"/path2/resource"`))
					Expect(journal).To(ContainSubstring(`"scheme":"http"`))
					Expect(journal).To(ContainSubstring(`"method":"POST"`))
					Expect(journal).To(ContainSubstring(`"body":"` + string(responseBody) + `"`))
					Expect(journal).To(ContainSubstring(`"status":200`))
				})
			})
		})
	})
})
