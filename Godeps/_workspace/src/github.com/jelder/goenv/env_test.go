package env_test

import (
	. "github.com/jelder/env"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestEnv(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Env Suite")
}

var ENV EnvMap

var _ = Describe("Book", func() {

	Context("With a malicious environment", func() {
		It("should refuse to load", func() {
			ENV, err := LoadEnvArrayString([]string{
				"OK=true",
				"Ok=false",
			})
			Expect(ENV).To(BeNil())
			Expect(err).To(HaveOccurred())
		})
	})

	Context("With a typical environment", func() {
		BeforeEach(func() {
			ENV = MustLoadEnvArrayString([]string{
				"KEY1=val1",
				"KEY2=val2",
				"BOOL1=true",
				"BOOL2=f",
				"FLOAT1=-3.14",
			})
		})
		It("Should work like a typical map", func() {
			Expect(ENV["KEY1"]).To(Equal("val1"))
		})
		It("Should allow fetching with a default", func() {
			Expect(ENV.Get("KEY1", "string")).To(Equal("val1"))
			Expect(ENV.Get("KEY3", "string")).To(Equal("string"))
		})

		Describe("IsSet", func() {
			BeforeEach(func() {
				ENV = MustLoadEnvArrayString([]string{
					"SET=lksdjf",
				})
			})
			It("Should return only presense", func() {
				Expect(ENV.IsSet("SET")).To(BeTrue())
				Expect(ENV.IsSet("NOPE")).To(BeFalse())
			})
		})

		Describe("GetBool", func() {
			It("Should treat any truthy value as true", func() {
				Expect(ENV.GetBool("BOOL1")).To(BeTrue())
			})
			It("Should treat any non-truthy values as false", func() {
				Expect(ENV.GetBool("BOOL2")).To(BeFalse())
				Expect(ENV.GetBool("BOOL99")).To(BeFalse())
			})
		})

		Describe("GetNumber", func() {
			Context("with a missing value", func() {
				It("should use the given default value", func() {
					Expect(ENV.GetNumber("FLOAT2", 123)).To(Equal(float64(123)))
				})
			})
			Context("with numbery value", func() {
				It("should parse and return the value", func() {
					Expect(ENV.GetNumber("FLOAT1", 123)).To(Equal(-3.14))
				})
			})
		})

	})
})
