using System.Collections.Specialized;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Net;
using System.Net.Sockets;
using System.Threading;
using System.Threading.Tasks;
using famous.oauth.requests;
using famous.oauth.responses;

namespace famous.oauth
{
  /// <summary>
  /// OAuth 2.0 verification code receiver that runs a local server on a free port and waits for a call with the 
  /// authorization verification code.
  /// </summary>
  internal class LocalServerCodeReceiver : ICodeReceiver
  {

    /// <summary>The call back format. Expects one port parameter.</summary>
    private const string LoopbackCallback = "http://localhost:{0}/authorize/";

    /// <summary>Close HTML tag to return the browser so it will close itself.</summary>
    private const string ClosePageResponse =
      @"<html>
  <head><title>OAuth 2.0 Authentication Token Received</title></head>
  <body>
    Received verification code. Closing...
    <script type='text/javascript'>
      window.setTimeout(function() {
          window.open('', '_self', ''); 
          window.close(); 
        }, 1000);
      if (window.opener) { window.opener.checkToken(); }
    </script>
  </body>
</html>";

    private string redirect_uri;
    public string CallbackUrl
    {
      get
      {
        if (!string.IsNullOrEmpty(redirect_uri))
        {
          return redirect_uri;
        }

        return redirect_uri = string.Format(LoopbackCallback, GetRandomUnusedPort());
      }
    }

    public async Task<AuthorizationCodeResponse> ReceiveCodeAsync(string authorizationUrl,
      CancellationToken taskCancellationToken)
    {
      using (var listener = new HttpListener())
      {
        listener.Prefixes.Add(CallbackUrl);
        try
        {
          listener.Start();

          var p = Process.Start(authorizationUrl);


          // Wait to get the authorization code response.
          var context = await listener.GetContextAsync().ConfigureAwait(false);
          var coll = context.Request.QueryString;

          // Write a "close" response.
          using (var writer = new StreamWriter(context.Response.OutputStream))
          {
            writer.WriteLine(ClosePageResponse);
            writer.Flush();
          }
          context.Response.OutputStream.Close();
          // Create a new response URL with a dictionary that contains all the response query parameters.
          return new AuthorizationCodeResponse(coll.AllKeys.ToDictionary(k => k, k => coll[(string) k]));
        }
        finally
        {
          listener.Close();
        }
      }
    }


    /// <summary>Returns a random, unused port.</summary>
    private static int GetRandomUnusedPort()
    {
      var listener = new TcpListener(IPAddress.Loopback, 0);
      try
      {
        listener.Start();
        return ((IPEndPoint)listener.LocalEndpoint).Port;
      }
      finally
      {
        listener.Stop();
      }
    }
  }
}