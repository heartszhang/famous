using System;
using famous.oauth.core;
using famous.oauth.core.request;
using famous.oauth.utils;

namespace famous.oauth.requests
{
    /// <summary>
    /// OAuth 2.0 request URL for an authorization web page to allow the end user to authorize the application to 
    /// access their protected resources and that returns an authorization code, as specified in 
    /// http://tools.ietf.org/html/rfc6749#section-4.1.
    /// </summary>
  internal class AuthorizationCodeRequest 
    {
        /// <summary>
        /// Gets or sets the response type which must be <c>code</c> for requesting an authorization code or 
        /// <c>token</c> for requesting an access token (implicit grant), or space separated registered extension 
        /// values. See http://tools.ietf.org/html/rfc6749#section-3.1.1 for more details
        /// </summary>
        [HttpRequestParameter("response_type", HttpRequestParameter.ParamType.Query)]
        public string ResponseType { get; set; }

        /// <summary>Gets or sets the client identifier.</summary>
        [HttpRequestParameter("client_id", HttpRequestParameter.ParamType.Query)]
        public string ClientId { get; set; }

        /// <summary>
        /// Gets or sets the URI that the authorization server directs the resource owner's user-agent back to the 
        /// client after a successful authorization grant, as specified in 
        /// http://tools.ietf.org/html/rfc6749#section-3.1.2 or <c>null</c> for none.
        /// </summary>
        [HttpRequestParameter("redirect_uri", HttpRequestParameter.ParamType.Query)]
        public string RedirectUri { get; set; }

        /// <summary>
        /// Gets or sets space-separated list of scopes, as specified in http://tools.ietf.org/html/rfc6749#section-3.3
        /// or <c>null</c> for none.
        /// </summary>
        [HttpRequestParameter("scope", HttpRequestParameter.ParamType.Query)]
        public string Scope { get; set; }

        /// <summary>
        /// Gets or sets the state (an opaque value used by the client to maintain state between the request and 
        /// callback, as mentioned in http://tools.ietf.org/html/rfc6749#section-3.1.2.2 or <c>null</c> for none.
        /// </summary>
        [HttpRequestParameter("state", HttpRequestParameter.ParamType.Query)]
        public string State { get; set; }

        /// <summary>
        /// Constructs a new authorization code request with the specified URI and sets response_type to <c>code</c>.
        /// </summary>
        public AuthorizationCodeRequest()
        {
            ResponseType = "code";
        }

        /// <summary>Creates a <seealso cref="System.Uri"/> which is used to request the authorization code.</summary>
        public Uri Build(Uri serverurl)
        {
          var b = new RequestBuilder(this){ BaseUri = serverurl};
 //         b.AddParameter(HttpRequestParameter.ParamType.Query, "x_required_offers", "Bing/Search");
          return b.BuildUri();
        }
    }
}
