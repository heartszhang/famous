using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Linq;
using System.Linq.Expressions;
using System.Reflection;
using System.Text;
using System.Threading.Tasks;

namespace famousfront.utils
{
    public static class NotificationExtensions
    {
        public static void Notify(this PropertyChangedEventHandler eventHandler, Expression<Func<object>> expression)
        {
            if (null == eventHandler)
            {
                return;
            }
            LambdaExpression lambda = expression;
            NotifyInternal(eventHandler, lambda);
        }

        public static void Notify<T>(this PropertyChangedEventHandler eventHandler, Expression<Func<T>> expression)
        {
            if (null == eventHandler)
            {
                return;
            }
            LambdaExpression lambda = expression;
            NotifyInternal(eventHandler, lambda);
        }

        private static void NotifyInternal(PropertyChangedEventHandler eventHandler, LambdaExpression lambda)
        {
            MemberExpression memberExpression;
            if (lambda.Body is UnaryExpression)
            {
                var unaryExpression = lambda.Body as UnaryExpression;
                memberExpression = unaryExpression.Operand as MemberExpression;
            }
            else
            {
                memberExpression = lambda.Body as MemberExpression;
            }
            var constantExpression = memberExpression.Expression as ConstantExpression;
            var propertyInfo = memberExpression.Member as PropertyInfo;

            eventHandler(constantExpression.Value, new PropertyChangedEventArgs(propertyInfo.Name));
        }
    }
}
