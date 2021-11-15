use std::io::Result;

fn main() -> Result<()> {
	prost_build::compile_protos(
		&["../proto/v1/g_counter.proto"],
		&["../proto/v1"]
	)?;
	Ok(())
}
