package garden_acceptance_test

import (
	"github.com/cloudfoundry-incubator/garden"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("bind_mounts", func() {
	It("can mount a read-only BindMount (#75464648)", func() {
		container := createContainer(gardenClient, garden.ContainerSpec{
			BindMounts: []garden.BindMount{
				garden.BindMount{
					SrcPath: "/var/vcap/packages",
					DstPath: "/home/alice/readonly",
					Mode:    garden.BindMountModeRO,
				},
			},
		})

		buffer := gbytes.NewBuffer()
		process, err := container.Run(garden.ProcessSpec{
			User: "root",
			Path: "ls",
			Args: []string{"/home/alice/readonly"},
		}, recordedProcessIO(buffer))
		Ω(err).ShouldNot(HaveOccurred())
		Ω(process.Wait()).Should(Equal(0))
		Ω(buffer).Should(gbytes.Say("rootfs"))

		process, err = container.Run(garden.ProcessSpec{
			User: "root",
			Path: "touch",
			Args: []string{"/home/alice/readonly/new_file"},
		}, silentProcessIO)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(process.Wait()).ShouldNot(Equal(0))

		process, err = container.Run(garden.ProcessSpec{
			User: "root",
			Path: "ls",
			Args: []string{"/home/alice/readonly"},
		}, recordedProcessIO(buffer))
		Ω(err).ShouldNot(HaveOccurred())
		Ω(process.Wait()).Should(Equal(0))
		Ω(buffer).ShouldNot(gbytes.Say("new_file"))
	})

	It("can mount a read/write BindMount (#75464648)", func() {
		container := createContainer(gardenClient, garden.ContainerSpec{
			BindMounts: []garden.BindMount{
				garden.BindMount{
					SrcPath: "/tmp",
					DstPath: "/home/alice/readwrite",
					Mode:    garden.BindMountModeRW,
					Origin:  garden.BindMountOriginContainer,
				},
			},
		})

		buffer := gbytes.NewBuffer()
		process, err := container.Run(garden.ProcessSpec{
			User: "root",
			Path: "ls",
			Args: []string{"/home/alice/readwrite"},
		}, recordedProcessIO(buffer))
		Ω(err).ShouldNot(HaveOccurred())
		Ω(process.Wait()).Should(Equal(0))
		Ω(buffer).ShouldNot(gbytes.Say("new_file"))

		buffer = gbytes.NewBuffer()
		process, err = container.Run(garden.ProcessSpec{
			User: "root",
			Path: "touch",
			Args: []string{"/home/alice/readwrite/new_file"},
		}, recordedProcessIO(buffer))
		Ω(err).ShouldNot(HaveOccurred())
		Ω(process.Wait()).Should(Equal(0))

		buffer = gbytes.NewBuffer()
		process, err = container.Run(garden.ProcessSpec{
			User: "root",
			Path: "ls",
			Args: []string{"/home/alice/readwrite"},
		}, recordedProcessIO(buffer))
		Ω(err).ShouldNot(HaveOccurred())
		Ω(process.Wait()).Should(Equal(0))
		Ω(buffer).Should(gbytes.Say("new_file"))
	})
})
