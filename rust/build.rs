use std::io::Result;

fn main() -> Result<()> {
    let read_dir = std::fs::read_dir("../proto/v1")?;
    let paths: Vec<_> = read_dir.map(Result::unwrap).map(|d| d.path()).collect();

    println!("{:?}", paths);
    prost_build::compile_protos(&paths, &["../proto/v1"])?;
    Ok(())
}
