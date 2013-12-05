using System;
using System.Linq;
using System.Reflection;
using famous.oauth.core.request;

namespace famous.oauth.utils
{
  /// <summary>
  /// Utility class for iterating on <seealso cref="HttpRequestParameter"/> properties in a request object.
  /// </summary>
  internal static class ParameterUtils
  {
    /// <summary>
    /// Iterates over all <seealso cref="HttpRequestParameter"/> properties in the request object and invokes
    /// the specified action for each of them.
    /// </summary>
    /// <param name="request">A request object</param>
    /// <param name="action">An action to invoke which gets the parameter type, name and its value</param>
    public static void IterateParameters(object request, Action<HttpRequestParameter.ParamType, string, object> action)
    {
      // Use reflection to build the parameter dictionary.
      foreach (var property in request.GetType().GetProperties(BindingFlags.Instance | BindingFlags.Public))
      {
        // Retrieve the RequestParameterAttribute.
        var attribute =
          property.GetCustomAttributes(typeof(HttpRequestParameter), false).FirstOrDefault() as HttpRequestParameter;
        if (attribute == null)
        {
          continue;
        }

        // Get the name of this parameter from the attribute, if it doesn't exist take a lower-case variant of 
        // property name.
        var name = attribute.Name ?? property.Name.ToLower();

        var propertyType = property.PropertyType;
        var value = property.GetValue(request, null);

        // Call action with the type name and value.
        if (propertyType.IsValueType || value != null)
        {
          action(attribute.Type, name, value);
        }
      }
    }

  }
}