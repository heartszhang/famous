using System;

namespace famous.oauth.responses
{
    /// <summary>
    /// Token response exception which is thrown in case of receiving a token error when an authorization code or an 
    /// access token is expected.
    /// </summary>
  public sealed class ResponseException<TResponse> : Exception
    {
        /// <summary>The error information.</summary>
    public TResponse Error { get; private set; }

        /// <summary>Constructs a new token response exception from the given error.</summary>
    public ResponseException(TResponse error)
            : base(error.ToString())
        {
            Error = error;
        }
    }
}