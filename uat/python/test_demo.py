import subprocess

def test_demo_binary():
    result = subprocess.run(["go", "run", "../../cmd/demo"], capture_output=True, text=True)
    assert result.returncode == 0, result.stderr
    assert result.stdout, "expected some output"
