use hyper::{Body, Method, Request, Response, Server, StatusCode};
use hyper::service::{service_fn, make_service_fn};
use std::error::Error;
use tokio::signal;

use skywalking::{reporter::grpc::GrpcReporter, trace::tracer::Tracer};


async fn handle_request(tracer: Tracer, req: Request<Body>) -> Result<Response<Body>, hyper::Error> {
    let mut ctx = tracer.create_trace_context();

    {
        let _span = ctx.create_entry_span("/root");

        let path = req.uri().path();
        let method = req.method();
        let mut response = Response::new(Body::empty());
        {
            if method == Method::GET && path == "/hello" {
                let _span2 = ctx.create_local_span("/hello");

                *response.body_mut() = Body::from("Hello, World!");

                Ok(response)
            } else if method == Method::POST && path == "/echo" {
                let _span3 = ctx.create_local_span("/echo");

                let whole_body = hyper::body::to_bytes(req.into_body()).await?;
                let body_str = String::from_utf8_lossy(&whole_body).to_string();
                *response.body_mut() = Body::from(body_str);

                Ok(response)
            } else {
                *response.status_mut() = StatusCode::NOT_FOUND;
                *response.body_mut() = Body::from("404 Not Found");

                Ok(response)
            }
        }
    }

}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    // Connect to skywalking oap server.
    let collector_backend_address = "<EndPoint>";
    let auth_token = "<TOKEN>";

    let service_name = "<ServiceName>";
    let listen_addr = ([0, 0, 0, 0], 9999).into();

    let reporter = GrpcReporter::connect(collector_backend_address).await?;
    // Optional authentication, based on backend setting.
    let reporter = reporter.with_authentication(auth_token);
    // Spawn the reporting in background, with listening the graceful shutdown signal.
    let handle = reporter
        .reporting()
        .await
        .with_graceful_shutdown(async move {
            signal::ctrl_c().await.expect("failed to listen for event");
        })
        .spawn();

    // Do tracing.
    let tracer = Tracer::new(service_name, "instance", reporter.clone());

    // 启动 HTTP 服务器，并等待请求的到来
    // let addr = ([0, 0, 0, 0], 9999).into();
    let server = Server::bind(&listen_addr).serve(make_service_fn(|_| {
        let tracer = tracer.clone();
        async {
            Ok::<_, hyper::Error>(service_fn(move |req|{
                let tracer = tracer.clone();
                // let cloned_tracer = cloned_tracer.clone();
                async move {
                    handle_request(tracer, req).await
                }
            }))
        }}));

    println!("Listening on http://{}", listen_addr);

    // 等待 HTTP 服务器退出
    server.await?;

    handle.await?;

    Ok(())
}
