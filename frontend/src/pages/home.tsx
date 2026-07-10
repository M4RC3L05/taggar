import { useLocation } from "wouter-preact";

const Home = () => {
  const [_, navigate] = useLocation();

  return (
    <main style={{ display: "flex", textAlign: "center", width: "100%", height: "100%", flexDirection: "column", justifyContent: "center", alignItems: "center" }}>
      <h1>Welcome to Tagga</h1>
      <p>Select a single file or a directory to get started</p>

      <div class="d-flex">
        <button
          type="button"
          style={{ marginRight: "8px" }}
          class="btn btn-primary"
          onClick={() => {
            globalThis.chooseFile()
              .then((x) => {
                if (x !== "") {
                  navigate(`/manage-file/${encodeURIComponent(x)}`);
                }
              })
              .catch((e) => console.error("open file err", e));
          }}
        >
          Select file
        </button>
        <button
          type="button"
          class="btn btn-primary"
          onClick={() => {
            globalThis.chooseDirectory()
              .then((x) => console.log("open directory", x))
              .catch((e) => console.error("open directory err", e));
          }}
        >
          Select folder
        </button>
      </div>
    </main>
  );
};

export default Home;
