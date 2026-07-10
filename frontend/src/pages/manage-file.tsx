import { FunctionComponent } from "preact";
import { useEffect, useState } from "preact/hooks";
import { Link, useParams } from "wouter-preact";

const ViewMetadata: FunctionComponent<{ file: string }> = ({ file }) => {
  const [metadata, setMetadata] = useState(null);

  useEffect(() => {
    globalThis.getMediaTags(file)
      .then((x) => {
        setMetadata(x)
      })
      .catch((e) => console.log("err", e));
  }, [file]);

  return (
    <div style={{ overflowY: "auto", paddingTop: "8px", paddingBottom: "8px" }} class="row">
      <div class="col-sm-4" style={{ display: "flex", alignItems: "center", flexDirection: "column" }}>
        <img
          style={{ aspectRatio: "1/1", maxWidth: "500px", width: "100%", height: "auto", marginBottom: "8px" }}
          src={metadata?.cover ??
            "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII="}
        />
      </div>
      <div class="col-sm-8">
        <p>
          <strong>Album Artist:</strong>{" "}
          <span>
            {metadata?.albumArtist ?? "-"}
          </span>
        </p>
        <p>
          <strong>Album:</strong>{" "}
          <span>
            {metadata?.album ?? "-"}
          </span>
        </p>
        <p>
          <strong>Title:</strong>{" "}
          <span>
            {metadata?.title ?? "-"}
          </span>
        </p>
        <p>
          <strong>Year:</strong>{" "}
          <span>
            {metadata?.year ?? "-"}
          </span>
        </p>
        <p>
          <strong>Artist:</strong>{" "}
          <span>
            {metadata?.artist ?? "-"}
          </span>
        </p>
        <p>
          <strong>Genre:</strong>{" "}
          <span>
            {metadata?.genre ?? "-"}
          </span>
        </p>
        <hr />
        <p>
          <strong>Track:</strong>{" "}
          <span>
            {metadata?.track ?? "-"}
          </span>
        </p>
        <p>
          <strong>Track Count:</strong>{" "}
          <span>
            {metadata?.trackCount ?? "-"}
          </span>
        </p>
        <p>
          <strong>Disc:</strong>{" "}
          <span>
            {metadata?.disc ?? "-"}
          </span>
        </p>
        <p>
          <strong>Disc Count:</strong>{" "}
          <span>
            {metadata?.discCount ?? "-"}
          </span>
        </p>
      </div>
    </div>
  );
};

const ManageFile: FunctionComponent = () => {
  const params = useParams<{ file: string }>();
  const selectedFile = decodeURIComponent(params.file);
  const [mode, setMode] = useState<"view" | "edit">("view");

  return (
    <div style={{ height: "100%", display: "flex", flexDirection: "column" }} class="container-fluid">
      <div style={{ flexShrink: "0", borderBottom: "1px solid var(--mu-tertiary)" }} class="row">
        <div style={{ display: "flex", flexDirection: "row", paddingTop: "5px", paddingBottom: "5px", alignItems: "center" }} class="col-12 border-bottom">
          <h3 style={{ margin: 0, flexGrow: "1" }}>
            <Link href="/" style={{ marginRight: "5px", textDecoration: "none" }}>
              <i class="bi bi-arrow-left" />
            </Link>
            Manage file {selectedFile}
          </h3>
          <div style={{ display: "flex" }}>
            <button
              style={{ marginBottom: 0 }}
              type="button"
              class={`me-2 btn ${mode === "view" ? "btn-primary" : "btn-secondary"
                }`}
              onClick={() => setMode("view")}
            >
              View
            </button>
            <button
              style={{ marginBottom: 0 }}
              type="button"
              class={`btn ${mode === "edit" ? "btn-primary" : "btn-secondary"
                }`}
              onClick={() => setMode("edit")}
            >
              Edit
            </button>
            {mode === "edit" && (
              <>
                <span class="mx-2 h-100 border-end" />
                <button style={{ marginBottom: 0 }} type="button" class="btn btn-success">Save</button>
                <button style={{ marginBottom: 0 }} type="button" class="btn">Reset</button>
              </>
            )}
          </div>
        </div>
      </div>
      {mode === "view" && <ViewMetadata file={selectedFile} />}
    </div>
  );
};

export default ManageFile;
