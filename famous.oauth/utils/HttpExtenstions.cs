using System.Net;
using System.Net.Http;

namespace famous.oauth.utils
{
    /// <summary>
    /// Extension methods to <see cref="System.Net.Http.HttpRequestMessage"/> and 
    /// <see cref="System.Net.Http.HttpResponseMessage"/>.
    /// </summary>
    static class HttpExtenstions
    {
        /// <summary>Returns <c>true</c> if the response contains one of the redirect status codes.</summary>
        public static bool IsRedirectStatusCode(this HttpResponseMessage message)
        {
            switch (message.StatusCode)
            {
                case HttpStatusCode.Moved:
                case HttpStatusCode.Redirect:
                case HttpStatusCode.RedirectMethod:
                case HttpStatusCode.TemporaryRedirect:
                    return true;
                default:
                    return false;
            }
        }

        /// <summary>Sets an empty HTTP content.</summary>
        public static HttpContent SetEmptyContent(this HttpRequestMessage request)
        {
            request.Content = new ByteArrayContent(new byte[0]);
            request.Content.Headers.ContentLength = 0;
            return request.Content;
        }
    }
}
