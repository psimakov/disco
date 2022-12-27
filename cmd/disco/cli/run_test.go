package cli

import (
	"context"
	"flag"
	"testing"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestImgCmd(t *testing.T) {
	imgCmd = testImgCmdFunc

	set := flag.NewFlagSet("", flag.ContinueOnError)
	set.String("project", "test", "test")

	c := cli.NewContext(app, set, nil)
	err := runImagesCmd(c)
	assert.NoError(t, err)
}

func testImgCmdFunc(ctx context.Context, in *types.ImagesQuery) error {
	return nil
}

func TestVulCmd(t *testing.T) {
	vulCmd = testVulCmdFunc

	set := flag.NewFlagSet("", flag.ContinueOnError)
	set.String("project", "test", "test")

	c := cli.NewContext(app, set, nil)
	err := runVulnsCmd(c)
	assert.NoError(t, err)
}

func testVulCmdFunc(ctx context.Context, in *types.VulnsQuery) error {
	return nil
}

func TestLicCmd(t *testing.T) {
	licCmd = testLicCmdFunc

	set := flag.NewFlagSet("", flag.ContinueOnError)
	set.String("project", "test", "test")

	c := cli.NewContext(app, set, nil)
	err := runLicenseCmd(c)
	assert.NoError(t, err)
}

func testLicCmdFunc(ctx context.Context, in *types.SimpleQuery) error {
	return nil
}