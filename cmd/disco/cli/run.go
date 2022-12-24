package cli

import (
	"github.com/mchmarny/disco/pkg/disco"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	c "github.com/urfave/cli/v2"
)

var (
	projectIDFlag = &c.StringFlag{
		Name:     "project",
		Aliases:  []string{"p"},
		Usage:    "project ID",
		Required: false,
	}

	outputPathFlag = &c.StringFlag{
		Name:     "output",
		Aliases:  []string{"o"},
		Usage:    "path where to save the output",
		Required: false,
	}

	outputFormatFlag = &c.StringFlag{
		Name:     "format",
		Aliases:  []string{"f"},
		Usage:    "output format (json, yaml, raw)",
		Required: false,
	}

	outputDigestOnlyFlag = &c.BoolFlag{
		Name:  "digest",
		Usage: "output only image digests",
		Value: false,
	}

	cveFlag = &c.StringFlag{
		Name:     "cve",
		Aliases:  []string{"e"},
		Usage:    "exposure ID (CVE number, e.g. CVE-2019-19378)",
		Required: false,
	}

	runCmd = &c.Command{
		Name:  "run",
		Usage: "Cloud Run commands",
		Subcommands: []*c.Command{
			{
				Name:    "images",
				Aliases: []string{"img", "i"},
				Usage:   "List deployed container images",
				Action:  runImagesCmd,
				Flags: []c.Flag{
					projectIDFlag,
					outputPathFlag,
					outputFormatFlag,
					outputDigestOnlyFlag,
				},
			},
			{
				Name:    "vulnerabilities",
				Aliases: []string{"vul", "v"},
				Usage:   "Check for OS-level exposures in deployed images (supports specific CVE filter)",
				Action:  runVulnsCmd,
				Flags: []c.Flag{
					projectIDFlag,
					outputPathFlag,
					outputFormatFlag,
					cveFlag,
				},
			},
		},
	}
)

func printVersion(c *c.Context) {
	log.Info().Msgf(c.App.Version)
}

func runImagesCmd(c *c.Context) error {
	fmtStr := c.String(outputFormatFlag.Name)
	outFmt, err := disco.ParseOutputFormat(fmtStr)
	if err != nil {
		return errors.Wrapf(err, "error parsing output format: %s", fmtStr)
	}

	in := &disco.ImagesQuery{}
	in.ProjectID = c.String(projectIDFlag.Name)
	in.OutputPath = c.String(outputPathFlag.Name)
	in.OutputFmt = outFmt
	in.OnlyDigest = c.Bool(outputDigestOnlyFlag.Name)

	printVersion(c)
	if err := disco.DiscoverImages(c.Context, in); err != nil {
		return errors.Wrap(err, "error discovering images")
	}

	return nil
}

func runVulnsCmd(c *c.Context) error {
	fmtStr := c.String(outputFormatFlag.Name)
	outFmt, err := disco.ParseOutputFormat(fmtStr)
	if err != nil {
		return errors.Wrapf(err, "error parsing output format: %s", fmtStr)
	}

	in := &disco.VulnsQuery{}
	in.ProjectID = c.String(projectIDFlag.Name)
	in.OutputPath = c.String(outputPathFlag.Name)
	in.CVE = c.String(cveFlag.Name)
	in.OutputFmt = outFmt

	printVersion(c)

	if in.CVE != "" {
		log.Info().Msg("Note: vulnerability scans currently limited to base OS only")
	}

	if err := disco.DiscoverVulns(c.Context, in); err != nil {
		return errors.Wrap(err, "error excuting command")
	}

	return nil
}