window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    url: "/swagger/swagger.json",
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout",
    requestInterceptor: (req) => {
      const token = localStorage.getItem('jwt-token'); // Храните JWT в localStorage
      if (token) {
        req.headers['Authorization'] = `Bearer ${token}`;
      }
      return req;
    }
  });

  //</editor-fold>
};
