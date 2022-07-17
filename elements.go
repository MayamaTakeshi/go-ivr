package main

import (
	"errors"
	"strings"

	"github.com/beevik/etree"
)

type attributeSpec struct {
	required     bool
	defaultValue string
	validator    func(string) (bool, error)
}

type elementSpec struct {
	attributes      map[string]attributeSpec
	allowedChildren []string
}

var elementSpecs map[string]elementSpec

func init() {
	elementSpecs = map[string]elementSpec{
		"Play": {
			attributes: map[string]attributeSpec{
				"loop": {
					required:     false,
					defaultValue: "1",
					validator:    func(value string) (bool, error) { return true, nil },
				},
				"terminators": {
					required:     false,
					defaultValue: "",
					validator:    func(value string) (bool, error) { return true, nil },
				},
				"purpose": {
					required:     false,
					defaultValue: "prompt",
					validator:    func(value string) (bool, error) { return true, nil },
				},
			},
			allowedChildren: []string{},
		},

		"Speak": {
			attributes: map[string]attributeSpec{
				"language": {
					required:     true,
					defaultValue: "",
					validator:    func(value string) (bool, error) { return true, nil },
				},
				"voice": {
					required:     true,
					defaultValue: "",
					validator:    func(value string) (bool, error) { return true, nil },
				},
				"loop": {
					required:     false,
					defaultValue: "1",
					validator:    func(value string) (bool, error) { return true, nil },
				},
				"purpose": {
					required:     false,
					defaultValue: "prompt",
					validator:    func(value string) (bool, error) { return true, nil },
				},
			},
			allowedChildren: []string{},
		},

		"Wait": {
			attributes: map[string]attributeSpec{
				"length": {
					required:     true,
					defaultValue: "1",
					validator:    func(value string) (bool, error) { return true, nil },
				},
				"purpose": {
					required:     false,
					defaultValue: "prompt",
					validator:    func(value string) (bool, error) { return true, nil },
				},
			},
			allowedChildren: []string{},
		},

		"GetDigits": {
			attributes: map[string]attributeSpec{
				"action": {
					required:     false,
					defaultValue: "",
					validator:    func(value string) (bool, error) { return true, nil },
				},
				"timeout": {
					required:     false,
					defaultValue: "5",
					validator:    func(value string) (bool, error) { return true, nil },
				},
				"interdigitTimeout": {
					required:     false,
					defaultValue: "5",
					validator:    func(value string) (bool, error) { return true, nil },
				},
				"finishOnKey": {
					required:     false,
					defaultValue: "#",
					validator:    func(value string) (bool, error) { return true, nil },
				},
				"numDigits": {
					required:     false,
					defaultValue: "99",
					validator:    func(value string) (bool, error) { return true, nil },
				},
				"maxTries": {
					required:     false,
					defaultValue: "1",
					validator:    func(value string) (bool, error) { return true, nil },
				},
				"playBeep": {
					required:     false,
					defaultValue: "false",
					validator:    func(value string) (bool, error) { return true, nil },
				},
				"validDigits": {
					required:     false,
					defaultValue: "1234567890*#",
					validator:    func(value string) (bool, error) { return true, nil },
				},
			},
			allowedChildren: []string{"Play", "Speak", "Wait"},
		},
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func validate(root *etree.Element) (bool, error) {
	for _, child := range root.ChildElements() {
		if _, ok := elementSpecs[child.Tag]; !ok {
			words := []string{"Unknown Element", child.Tag}
			return false, errors.New(strings.Join(words, " "))
		}

		for _, attr := range child.Attr {
			if attrSpec, ok := elementSpecs[child.Tag].attributes[attr.Key]; !ok {
				words := []string{"Element", child.Tag, "unknown Attribute", attr.Key}
				return false, errors.New(strings.Join(words, " "))
			} else {
				if ok, err := attrSpec.validator(attr.Value); !ok {
					return false, err
				}
			}
		}

		for _, childOfChild := range child.ChildElements() {
			if !contains(elementSpecs[child.Tag].allowedChildren, childOfChild.Tag) {
				words := []string{"Element", child.Tag, "cannot contain Element", childOfChild.Tag}
				return false, errors.New(strings.Join(words, " "))
			}
		}
		_, err := validate(child)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}