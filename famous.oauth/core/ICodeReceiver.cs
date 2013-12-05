using System.Threading;
using System.Threading.Tasks;
using famous.oauth.responses;

namespace famous.oauth.core
{
  /// <summary>OAuth 2.0 verification code receiver.</summary>
  public interface ICodeReceiver
  {
    /// <summary>Gets the redirected URI.</summary>
    string CallbackUrl { get; }

    /// <summary>Receives the authorization code.</summary>
    /// <param name="authorizationurl">The authorization code request URL</param>
    /// <param name="taskCancellationToken">Cancellation token</param>
    /// <returns>The authorization code response</returns>
    Task<AuthorizationCodeResponse> ReceiveCodeAsync(string authorizationurl,CancellationToken taskCancellationToken);
  }
}