#!/usr/bin/env python3
import os
import stat
import subprocess
import tempfile
import unittest
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[2]
GENERATE_SCRIPT = REPO_ROOT / "scripts" / "release" / "generate_artifacts.sh"
VERIFY_SCRIPT = REPO_ROOT / "scripts" / "release" / "verify_release_artifacts.sh"


def run_script(script: Path, env: dict[str, str]) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        ["bash", str(script)],
        cwd=str(REPO_ROOT),
        env=env,
        text=True,
        capture_output=True,
        check=False,
    )


class ReleaseArtifactsScriptTests(unittest.TestCase):
    def test_generate_artifacts_produces_checksum_and_signature_metadata(self) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            dist_dir = Path(temp_dir) / "dist"
            mock_bin_dir = Path(temp_dir) / "mock-bin"
            mock_bin_dir.mkdir(parents=True)

            fake_go = mock_bin_dir / "go"
            fake_go.write_text(
                "#!/usr/bin/env bash\n"
                "set -euo pipefail\n"
                'out=""\n'
                "for ((i=1; i <= $#; i++)); do\n"
                '  if [[ "${!i}" == "-o" ]]; then\n'
                "    j=$((i+1))\n"
                '    out="${!j}"\n'
                "  fi\n"
                "done\n"
                'if [[ -z "$out" ]]; then\n'
                "  echo 'missing -o output path' >&2\n"
                "  exit 1\n"
                "fi\n"
                'mkdir -p "$(dirname "$out")"\n'
                "printf 'fake weave binary' > \"$out\"\n",
                encoding="utf-8",
            )
            fake_go.chmod(fake_go.stat().st_mode | stat.S_IEXEC)

            env = os.environ.copy()
            env.update(
                {
                    "DIST_DIR": str(dist_dir),
                    "WEAVE_VERSION": "test",
                    "WEAVE_TARGETS": "linux/amd64",
                    "PATH": f"{mock_bin_dir}:{env.get('PATH', '')}",
                }
            )

            result = run_script(GENERATE_SCRIPT, env)
            self.assertEqual(result.returncode, 0, msg=result.stderr)
            self.assertIn("Artifacts generated in", result.stdout)

            tar_file = dist_dir / "weave_test_linux_amd64.tar.gz"
            checksums = dist_dir / "checksums.txt"
            signature = dist_dir / "checksums.txt.sig"
            pub_key = dist_dir / "checksums.pub"

            self.assertTrue(tar_file.exists(), "expected release artifact tarball")
            self.assertTrue(checksums.exists(), "expected checksums.txt")
            self.assertTrue(signature.exists(), "expected checksums.txt.sig")
            self.assertTrue(pub_key.exists(), "expected checksums.pub")
            self.assertIn(tar_file.name, checksums.read_text(encoding="utf-8"))

    def test_verify_release_artifacts_fails_when_signature_missing(self) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            dist_dir = Path(temp_dir) / "dist"
            dist_dir.mkdir(parents=True)
            (dist_dir / "checksums.txt").write_text("", encoding="utf-8")

            env = os.environ.copy()
            env["DIST_DIR"] = str(dist_dir)

            result = run_script(VERIFY_SCRIPT, env)
            self.assertNotEqual(result.returncode, 0)
            self.assertIn("missing checksums.txt.sig", result.stderr)

    def test_verify_release_artifacts_passes_with_generated_ephemeral_public_key(
        self,
    ) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            dist_dir = Path(temp_dir) / "dist"
            mock_bin_dir = Path(temp_dir) / "mock-bin"
            mock_bin_dir.mkdir(parents=True)

            fake_go = mock_bin_dir / "go"
            fake_go.write_text(
                "#!/usr/bin/env bash\n"
                "set -euo pipefail\n"
                'out=""\n'
                "for ((i=1; i <= $#; i++)); do\n"
                '  if [[ "${!i}" == "-o" ]]; then\n'
                "    j=$((i+1))\n"
                '    out="${!j}"\n'
                "  fi\n"
                "done\n"
                'mkdir -p "$(dirname "$out")"\n'
                "printf 'fake weave binary' > \"$out\"\n",
                encoding="utf-8",
            )
            fake_go.chmod(fake_go.stat().st_mode | stat.S_IEXEC)

            env = os.environ.copy()
            env.update(
                {
                    "DIST_DIR": str(dist_dir),
                    "WEAVE_VERSION": "verify",
                    "WEAVE_TARGETS": "linux/amd64",
                    "PATH": f"{mock_bin_dir}:{env.get('PATH', '')}",
                }
            )

            generate_result = run_script(GENERATE_SCRIPT, env)
            self.assertEqual(generate_result.returncode, 0, msg=generate_result.stderr)

            verify_result = run_script(VERIFY_SCRIPT, env)
            self.assertEqual(verify_result.returncode, 0, msg=verify_result.stderr)
            self.assertIn("Release artifact verification passed", verify_result.stdout)


if __name__ == "__main__":
    unittest.main()
